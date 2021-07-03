package linkedin

import (
	"time"

	"cloud.google.com/go/civil"
)

type AdLocale struct {
	Country  string `json:"country"`
	Language string `json:"language"`
}

type AdBudget struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

type AdTimestamp struct {
	Time int64 `json:"time"`
}

type AdVersion struct {
	VersionTag string `json:"versionTag"`
}

type AdChangeAuditStamps struct {
	Created      AdTimestamp `json:"created"`
	LastModified AdTimestamp `json:"lastModified"`
}

type AdRunSchedule struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

type AdDate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type AdDateRange struct {
	End   *AdDate `json:"end"`
	Start *AdDate `json:"start"`
}

type TimeGranularity string

const (
	TimeGranularityAll     TimeGranularity = "ALL"
	TimeGranularityDaily   TimeGranularity = "DAILY"
	TimeGranularityMonthly TimeGranularity = "MONTHLY"
	TimeGranularityYearly  TimeGranularity = "YEARLY"
)

func (d *AdDate) ToDate() *civil.Date {
	if d == nil {
		return nil
	}

	return &civil.Date{
		Year:  d.Year,
		Month: time.Month(d.Month),
		Day:   d.Day,
	}
}

func NewAdDate(date *civil.Date) *AdDate {
	if date == nil {
		return nil
	}

	return &AdDate{
		Year:  date.Year,
		Month: int(date.Month),
		Day:   date.Day,
	}
}
