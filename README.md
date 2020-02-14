AWS Cost (Golang)
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Codecov Coverage][11]][12] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22]

[1]: https://godoc.org/github.com/evalphobia/awscost?status.svg
[2]: https://godoc.org/github.com/evalphobia/awscost
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/awscost.svg
[6]: https://github.com/evalphobia/awscost/releases/latest
[7]: https://github.com/evalphobia/awscost/workflows/test/badge.svg
[8]: https://github.com/evalphobia/awscost/actions?query=workflow%3Atest
[9]: https://coveralls.io/repos/evalphobia/awscost/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/awscost?branch=master
[11]: https://codecov.io/github/evalphobia/awscost/coverage.svg?branch=master
[12]: https://codecov.io/github/evalphobia/awscost?branch=master
[13]: https://goreportcard.com/badge/github.com/evalphobia/awscost
[14]: https://goreportcard.com/report/github.com/evalphobia/awscost
[15]: https://img.shields.io/github/downloads/evalphobia/awscost/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/awscost/releases
[17]: https://img.shields.io/github/stars/evalphobia/awscost.svg
[18]: https://github.com/evalphobia/awscost/stargazers
[19]: https://codeclimate.com/github/evalphobia/awscost/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/awscost
[21]: https://bettercodehub.com/edge/badge/evalphobia/awscost?branch=master
[22]: https://bettercodehub.com/

`awscost` is library to get AWS estimate billing costs.
This library probides two ways to get costs.

- ClowdWatch API
- CostExplorer API

# Quick Usage


```go
import (
	"fmt"
	"time"

	"github.com/evalphobia/awscost/costexplorer"
)

func main() {
	endDate := time.Now().AddDate(0, 0, -1)

	// costs, err := cloudwatch.FetchDailyCost(endDate)
	costs, err := costexplorer.FetchDailyCost(endDate)
	if err != nil {
		return panic(err)
	}

	// set date to show it on the report
	costs.SetDate(endDate)
	text := costs.FormatAsOutputReport()

	// show cost data
	fmt.Println(text)

	// post cost data to slack channel
	err = sendCostToSlack(text)
	if err != nil {
		return panic(err)
	}
}
```



## Environment variables

|Name|Description|
|:--|:--|
| `AWSCOST_LANG` | Supported languages are `en`, `ja`.  [See the language file](/i18n_lang.go) |
