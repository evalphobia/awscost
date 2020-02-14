package awscost

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Costs struct {
	Total    float64
	Other    float64
	Services map[string]float64
	Date     string
}

func (c *Costs) SetDate(dt time.Time) {
	c.Date = dt.Format("2006-01-02")
}

func (c Costs) FormatAsOutputReport() string {
	svc := c.Services
	kvList := KVList{}
	for key, val := range svc {
		kvList = append(kvList, KV{
			Key:   key,
			Value: val,
		})
	}
	sort.Sort(kvList)

	results := make([]string, 0, len(svc)*2)
	results = append(results, Message("[AWS Estimate Costs] %s", c.Date))
	results = append(results, fmt.Sprintf("- Total:\t$%.2f", c.Total))
	results = append(results, "------------------------")
	for _, kv := range kvList {
		results = append(results, fmt.Sprintf("- %s:\t$%.2f", kv.Key, kv.Value))
	}
	results = append(results, fmt.Sprintf("- (Other):\t$%.2f", c.Other))
	return strings.Join(results, "\n")
}

type KV struct {
	Key   string
	Value float64
}

type KVList []KV

func (l KVList) Len() int {
	return len(l)
}

func (l KVList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l KVList) Less(i, j int) bool {
	if l[i].Value == l[j].Value {
		return l[i].Key < l[j].Key
	}
	return l[i].Value > l[j].Value
}
