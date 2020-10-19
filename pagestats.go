package linkedin

import "encoding/json"

type TotalPageStatistics struct {
	ClicksRaw json.RawMessage `json:"clicks"`
	Clicks    map[string]map[string]int64
	ViewsRaw  json.RawMessage `json:"views"`
	Views     map[string]PageViews
}

type TotalPageClicks struct {
	CareersPagePromoLinksClicks  int `json:"careersPagePromoLinksClicks"`
	CareersPageBannerPromoClicks int `json:"careersPageBannerPromoClicks"`
	CareersPageJobsClicks        int `json:"careersPageJobsClicks"`
	CareersPageEmployeesClicks   int `json:"careersPageEmployeesClicks"`
}

type PageViews struct {
	PageViews       int64 `json:"pageViews"`
	UniquePageViews int64 `json:"uniquePageViews"`
}
