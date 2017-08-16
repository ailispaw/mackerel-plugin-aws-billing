package mpawsbilling

import (
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type AWSBillingPlugin struct {
	Prefix  string
	Metrics []MetricValue
}

func (p AWSBillingPlugin) getPrefix() string {
	var prefix []string
	for i, metric := range p.Metrics {
		names := strings.Split(metric.Name, ".")
		names = names[:(len(names) - 1)]
		if i == 0 {
			prefix = names
			continue
		}
		for j, name := range names {
			if prefix[j] != name {
				prefix[j] = "#"
			}
		}
	}
	return strings.Join(prefix, ".")
}

func (p AWSBillingPlugin) MetricKeyPrefix() string {
	return "AWS-Billing"
}

func (p AWSBillingPlugin) GraphDefinition() map[string]mp.Graphs {
	metrics := func() []mp.Metrics {
		var metrics []mp.Metrics
		var labels []string
		for _, metric := range p.Metrics {
			names := strings.Split(metric.Name, ".")
			name := names[len(names)-1]
			if exists, _ := InArray(name, labels); !exists {
				metrics = append(metrics, mp.Metrics{
					Name:         name,
					Label:        name,
					AbsoluteName: true,
				})
				labels = append(labels, name)
			}
		}
		return metrics
	}()

	return map[string]mp.Graphs{
		p.Prefix: mp.Graphs{
			Label:   p.Prefix,
			Unit:    "float",
			Metrics: metrics,
		},
	}
}

func (p AWSBillingPlugin) FetchMetrics() (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	for _, metric := range p.Metrics {
		metrics[metric.Name] = metric.Value
	}

	return metrics, nil
}

func SendMetricsToMackerelHost(metrics []MetricValue) {
	var mpAWSBilling AWSBillingPlugin

	mpAWSBilling.Metrics = metrics
	mpAWSBilling.Prefix = mpAWSBilling.getPrefix()

	helper := mp.NewMackerelPlugin(mpAWSBilling)

	helper.Run()
}
