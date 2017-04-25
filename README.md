mackerel-plugin-aws-billing
=======================

# Overview
mackerel-plugin-aws-billing is mackerel agent plugin that gets AWS cost and makes a graph.

## Description

It uses [AWS CloudWatch Api](https://aws.amazon.com/ja/documentation/cloudwatch/) to get AWS Billing Data related to AWS Account.  
It can make a graph on [Service Metric](https://mackerel.io/ja/features/service-metrics/).  
It gets AWS cost every hour by using AWS API because unlike other metrics, AWS updates cost once in a few hours.  
It writes cache file on server every hour, and uses cache to output data.  

(※In order to output data to Service Metric, this plugin sends data to mackerel instead of mackerel-agent because mackerel-agent can not send data to Service Metric.)  

Caution: You must must enable Billing Alerts.(see https://docs.aws.amazon.com/ja_jp/awsaccountbilling/latest/aboutv2/monitor-charges.html)  

## Usage

You must set environment variables in `.env` file first.

- MACKEREL_METRICS(required)  
  Use ServiceMetric.

- MACKEREL_API_KEY(required)  
  MACKEREL_API_KEY is an API Key published by mackerel. 
  You must give read and write premissions to API Key.

- MACKEREL_SERVICE(required)  
  Specify mackerel service name to make graph on Service Metric page.

- AWS_ACCESS_KEY_ID(required)  
  AWS_ACCESS_KEY_ID is published by AWS. It is required to use AWS API.

- AWS_SECRET_ACCESS_KEY(required)  
  AWS_SECRET_ACCESS_KEY is published by AWS. It is required to use AWS API..

- AWS_TARGET_SERVICE(optional)  
  If AWS_TARGET_SERVICE is not specified, this plugin gets all available services.  
  If AWS_TARGET_SERVICE is specified, this plugin gets specified available services.  
  If you want to specify multiple services, separate with comma.(ex: AWS_TARGET_SERVICE=AmazonEC2,AWSLambda)  
  If you want to get sum of costs, give All to target.(ex: AWS_TARGET_SERVICE=All)

- AWS_DIMENSION_CURRENCY(optional)  
  Defalut is in US Dollar(USD).

```shell
$ docker run --name mackerel-plugin-aws-billing --env-file .env ailispaw/mackerel-plugin-aws-billing
$ echo "3 * * * * docker restart mackerel-plugin-aws-billing" | crontab -
```
