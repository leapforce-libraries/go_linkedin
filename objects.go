package linkedin

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

type Date struct {
	Day   int64 `json:"day"`
	Month int64 `json:"month"`
	Year  int64 `json:"year"`
}

type DateRange struct {
	End   Date `json:"end"`
	Start Date `json:"start"`
}

type TimeGranularity string

const (
	TimeGranularityAll     TimeGranularity = "ALL"
	TimeGranularityDaily   TimeGranularity = "DAILY"
	TimeGranularityMonthly TimeGranularity = "MONTHLY"
	TimeGranularityYearly  TimeGranularity = "YEARLY"
)
