package mpawsbilling

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type AWSBilling struct {
	Region      string
	Currency    string
	Credentials *credentials.Credentials
	CloudWatch  *cloudwatch.CloudWatch
}

type MetricValue struct {
	Name  string  `json:"name"`
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

func (b AWSBilling) GetServiceNameList() (targets []string) {
	out, err := b.CloudWatch.ListMetrics(
		&cloudwatch.ListMetricsInput{
			Namespace: aws.String("AWS/Billing"),
		})
	if err != nil {
		log.Fatalf("Failed to ListMetrics: %v", err)
	}

	for _, metric := range out.Metrics {
		for _, dimension := range metric.Dimensions {
			if *dimension.Name == "ServiceName" {
				targets = append(targets, *dimension.Value)
			}
		}
	}

	return targets
}

func (b AWSBilling) GetMetricValue(target string) (*MetricValue, error) {
	var dimensions = []*cloudwatch.Dimension{
		&cloudwatch.Dimension{
			Name:  aws.String("Currency"),
			Value: aws.String(b.Currency),
		},
		&cloudwatch.Dimension{
			Name:  aws.String("ServiceName"),
			Value: aws.String(target),
		},
	}

	now := time.Now()
	startTime := time.Unix(now.Unix()-(24*60*60), int64(now.Nanosecond()))

	out, err := b.CloudWatch.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: dimensions,
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(now),
		Namespace:  aws.String("AWS/Billing"),
		MetricName: aws.String("EstimatedCharges"),
		Period:     aws.Int64(60 * 60),
		Statistics: []*string{aws.String("Maximum")},
	})

	if err != nil {
		return nil, err
	}

	datapoints := out.Datapoints
	if len(datapoints) == 0 {
		return nil, errors.New("no datapoints")
	}

	var latest time.Time
	var latestIndex int

	for i, datapoint := range datapoints {
		if datapoint.Timestamp.After(latest) {
			latest = *datapoint.Timestamp
			latestIndex = i
		}
	}

	return &MetricValue{
		Name:  target,
		Time:  (*datapoints[latestIndex].Timestamp).Unix(),
		Value: *datapoints[latestIndex].Maximum,
	}, nil
}

func Do() {
	var (
		optDebug  bool
		optDryRun bool
		optHelp   bool
	)

	flag.BoolVar(&optDebug, "d", false, "Debug mode")
	flag.BoolVar(&optDryRun, "n", false, "Not send metrics to Mackerel, just check metrics with -d")
	flag.BoolVar(&optHelp, "h", false, "Show this help message")
	flag.Parse()

	if optHelp {
		flag.PrintDefaults()
		return
	}

	if os.Getenv("DEBUG") != "" {
		optDebug = true
	}

	optAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if optAccessKeyID == "" {
		log.Fatal("Please set AWS_ACCESS_KEY_ID environment variable.")
	}

	optSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if optSecretAccessKey == "" {
		log.Fatal("Please set AWS_SECRET_ACCESS_KEY environment variable.")
	}

	optCurrency := os.Getenv("AWS_DIMENSION_CURRENCY")
	if optCurrency == "" {
		optCurrency = "USD"
	}

	var billing AWSBilling

	billing.Region = "us-east-1"
	billing.Currency = optCurrency
	billing.Credentials = credentials.NewStaticCredentials(optAccessKeyID, optSecretAccessKey, "")

	billing.CloudWatch = cloudwatch.New(session.New(
		&aws.Config{
			Credentials: billing.Credentials,
			Region:      aws.String(billing.Region),
		}))

	var targets []string

	optTarget := os.Getenv("AWS_TARGET_SERVICE")
	if optTarget == "" {
		targets = billing.GetServiceNameList()
	} else {
		targets = strings.Split(optTarget, ",")
	}

	goroutines := len(targets)
	c := make(chan *MetricValue)
	for _, target := range targets {
		go func(m chan<- *MetricValue, target string) {
			metric, err := billing.GetMetricValue(target)
			if err == nil {
				m <- metric
			} else {
				log.Printf("%v: %s", target, err)
				m <- nil
			}

		}(c, target)
	}

	var metrics []MetricValue

	for i := 0; i < goroutines; i++ {
		metric := <-c
		if (metric != nil) && ((*metric).Value > 0) {
			metrics = append(metrics, MetricValue{
				Name:  "AWS.billing." + (*metric).Name,
				Time:  (*metric).Time,
				Value: (*metric).Value,
			})
		}
	}

	close(c)

	if optDebug {
		PrintInJSON(os.Stdout, metrics)
	}

	if len(metrics) > 0 {
		log.Printf("[AWS-Billing]: %v on %s", targets,
			time.Unix(metrics[0].Time, 0).Format(time.UnixDate))

		if optDryRun {
			return
		}

		optMackerelServiceName := os.Getenv("MACKEREL_SERVICE")
		if optMackerelServiceName != "" {
			optMackerelApiKey := os.Getenv("MACKEREL_API_KEY")

			if optMackerelApiKey == "" {
				log.Fatal("Please set MACKEREL_API_KEY environment variable.")
			}

			SendMetricsToMackerelService(optMackerelApiKey, optMackerelServiceName, metrics)
		} else {
			SendMetricsToMackerelHost(metrics)
		}
	}
}
