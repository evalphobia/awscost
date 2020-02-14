package cloudwatch

import (
	"fmt"
	"time"

	"github.com/evalphobia/aws-sdk-go-wrapper/cloudwatch"

	"github.com/evalphobia/awscost"
)

var emptyCost = awscost.Costs{}

func FetchDailyCost(endTime time.Time, targetServices ...string) (awscost.Costs, error) {
	endTime = autoFillEndTime(endTime)
	return Fetch(endTime.AddDate(0, 0, -1), endTime)
}

func Fetch(startTime, endTime time.Time, targetServices ...string) (awscost.Costs, error) {
	// Validation for AWS client
	// e.g.) AWS_ACCESS_KEY_ID is empty
	if err := validateCloudWatchClient(); err != nil {
		return emptyCost, fmt.Errorf("[validateCloudWatchClient]\t`%w`", err)
	}

	endTime = autoFillEndTime(endTime)
	targetServices = autoFillServices(targetServices)

	// Get costs of the target services.
	costs, err := fetchAllCosts(startTime, endTime, targetServices)
	if err != nil {
		return emptyCost, fmt.Errorf("[fetchAllCosts]\t`%w`", err)
	}
	return costs, nil
}

// Return date from the variable, or use yesterday.
func autoFillEndTime(endTime time.Time) time.Time {
	if !endTime.IsZero() {
		return endTime
	}
	return time.Now().In(time.UTC).AddDate(0, 0, -1)
}

// Return AWS services from the variable, or use default set.
func autoFillServices(services []string) []string {
	if len(services) != 0 {
		return services
	}
	return defaultAWSServices
}

func fetchAllCosts(startTime, endTime time.Time, targetServices []string) (awscost.Costs, error) {
	c := awscost.Costs{}
	total, err := fetchCosts(startTime, endTime)
	if err != nil {
		return c, err
	}

	svcTotal := 0.0
	svcCost := make(map[string]float64, len(targetServices))
	for _, s := range targetServices {
		cost, err := fetchCosts(startTime, endTime, s)
		if err != nil {
			return c, err
		}
		svcCost[s] = cost
		svcTotal += cost
	}
	c.Total = total
	c.Other = total - svcTotal
	c.Services = svcCost
	return c, nil
}

func fetchCosts(startTime, endTime time.Time, serviceName ...string) (float64, error) {
	resp1, err := fetchMetricsForCosts(startTime, endTime, serviceName...)
	if err != nil {
		return 0, err
	}
	cost := resp1.GetFirstValue()
	// 1st day does not need a diff of the previous day.
	if endTime.Day() == 1 {
		return cost, nil
	}

	resp2, err := fetchMetricsForCosts(startTime.AddDate(0, 0, -1), startTime, serviceName...)
	if err != nil {
		return 0, err
	}
	costPrevDay := resp2.GetFirstValue()
	return cost - costPrevDay, nil
}

func fetchMetricsForCosts(startTime, endTime time.Time, serviceName ...string) (Datapoints, error) {
	input := cloudwatch.MetricStatisticsInput{
		Namespace:  "AWS/Billing",
		MetricName: "EstimatedCharges",
		DimensionsMap: map[string]string{
			"Currency": "USD",
		},
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     int64(endTime.Sub(startTime).Seconds()),
		Statistics: []string{"Maximum"},
	}
	if len(serviceName) != 0 {
		input.DimensionsMap["ServiceName"] = serviceName[0]
	}

	return fetchCloudWatchMetrics(input)
}

var defaultAWSServices = []string{
	"AmazonApiGateway",
	"AmazonCloudWatch",
	"AmazonEC2",
	"AmazonECR",
	"AmazonDynamoDB",
	"AmazonElastiCache",
	"AmazonES",
	"AmazonGuardDuty",
	"AmazonInspector",
	"AmazonKinesis",
	"AmazonKinesisFirehose",
	"AmazonRDS",
	"AmazonRekognition",
	"AmazonRoute53",
	"AmazonS3",
	"AmazonSageMaker",
	"AmazonSES",
	"AmazonSNS",
	"AWSDataTransfer",
	"AWSIoT",
	"AWSLambda",
	"AWSQueueService",
	"CodeBuild",
}
