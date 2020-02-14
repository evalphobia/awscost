package costexplorer

import (
	"sync"
	"time"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/costexplorer"
)

var ceOnce sync.Once
var ceCli *costexplorer.CostExplorer

func getOrCreateCostExplorerClient() (*costexplorer.CostExplorer, error) {
	var err error
	ceOnce.Do(func() {
		ceCli, err = costexplorer.New(config.Config{})
	})
	return ceCli, err
}

func validateCostExplorerClient() error {
	_, err := getOrCreateCostExplorerClient()
	return err
}

func fetchCostData(startTime, endTime time.Time) (costexplorer.UsageResult, error) {
	cli, err := getOrCreateCostExplorerClient()
	if err != nil {
		return costexplorer.UsageResult{}, err
	}

	input := costexplorer.GetCostAndUsageInput{
		GroupByDimensionService: true,
		MetricUnblendedCost:     true,
		TimePeriodStart:         startTime,
		TimePeriodEnd:           endTime,
	}
	if startTime.AddDate(0, 0, 29).After(endTime) {
		input.GranularityMonthly = true
	}
	return cli.GetCostAndUsage(input)
}
