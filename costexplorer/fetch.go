package costexplorer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/evalphobia/awscost"
)

var emptyCost = awscost.Costs{}

func FetchDailyCost(endTime time.Time) (awscost.Costs, error) {
	endTime = autoFillEndTime(endTime)
	return Fetch(endTime.AddDate(0, 0, -1), endTime)
}

func Fetch(startTime, endTime time.Time) (awscost.Costs, error) {
	// Validation for AWS client
	// e.g.) AWS_ACCESS_KEY_ID is empty
	if err := validateCostExplorerClient(); err != nil {
		return emptyCost, fmt.Errorf("[validateCostExplorerClient]\t`%w`", err)
	}

	// Get costs of the target services.
	costs, err := fetchAllCosts(startTime, endTime, nil)
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

func fetchAllCosts(startTime, endTime time.Time, targetServices []string) (awscost.Costs, error) {
	c := awscost.Costs{}

	data, err := fetchCostData(startTime, endTime)
	if err != nil {
		return c, err
	}
	results := data.ResultsByTime
	if len(results) == 0 {
		return c, fmt.Errorf("[fetchCostData] result data is empty")
	}

	result := results[len(results)-1]
	groups := result.Groups
	if len(groups) == 0 {
		return c, fmt.Errorf("[fetchCostData] result group is empty")
	}

	svcTotal := 0.0
	svcCost := make(map[string]float64, len(groups))
	for _, g := range groups {
		if len(g.Keys) == 0 {
			continue
		}
		s := g.Keys[0]

		amount, _ := g.GetOne()
		cost, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			continue
		}

		svcCost[s] = cost
		svcTotal += cost
	}
	c.Total = svcTotal
	c.Services = svcCost
	return c, nil
}
