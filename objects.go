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
