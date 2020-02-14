package cloudwatch

import (
	"sync"
	"time"

	"github.com/evalphobia/aws-sdk-go-wrapper/cloudwatch"
	"github.com/evalphobia/aws-sdk-go-wrapper/config"
)

var cwOnce sync.Once
var cwCli *cloudwatch.CloudWatch

func getOrCreateCloudWatchClient() (*cloudwatch.CloudWatch, error) {
	var err error
	cwOnce.Do(func() {
		cwCli, err = cloudwatch.New(config.Config{})
	})
	return cwCli, err
}

func validateCloudWatchClient() error {
	_, err := getOrCreateCloudWatchClient()
	return err
}

func fetchCloudWatchMetrics(input cloudwatch.MetricStatisticsInput) (Datapoints, error) {
	cli, err := getOrCreateCloudWatchClient()
	if err != nil {
		return nil, err
	}
	resp, err := cli.GetMetricStatistics(input)
	if err != nil {
		return nil, err
	}

	return NewDatapoints(input, resp.Datapoints), nil
}

type Datapoint struct {
	MetricName string
	Value      float64
	Time       time.Time
}

type Datapoints []Datapoint

func NewDatapoints(input cloudwatch.MetricStatisticsInput, list []cloudwatch.Datapoint) Datapoints {
	if len(list) == 0 {
		return nil
	}

	name := input.MetricName
	data := make([]Datapoint, len(list))
	for i, p := range list {
		d := Datapoint{
			MetricName: name,
			Time:       p.Timestamp,
			Value:      p.Maximum,
		}
		data[i] = d
	}
	return data
}

func (p Datapoints) GetFirstValue() float64 {
	if len(p) == 0 {
		return 0
	}

	return p[0].Value
}
