package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type LifetimePageStatsResponse struct {
	Paging   Paging              `json:"paging"`
	Elements []LifetimePageStats `json:"elements"`
}

type LifetimePageStats struct {
	ByStaffCountRange []LifetimePageStatisticsByType `json:"pageStatisticsByStaffCountRange"`
	ByFunction        []LifetimePageStatisticsByType `json:"pageStatisticsByFunction"`
	BySeniority       []LifetimePageStatisticsByType `json:"pageStatisticsBySeniority"`
	ByIndustry        []LifetimePageStatisticsByType `json:"pageStatisticsByIndustry"`
	ByRegion          []LifetimePageStatisticsByType `json:"pageStatisticsByRegion"`
	ByCountry         []LifetimePageStatisticsByType `json:"pageStatisticsByCountry"`
	Totals            LifetimeTotalPageStatistics    `json:"totalPageStatistics"`
	Organization      string                         `json:"organization"`
}

type LifetimePageStatisticsByType struct {
	PageStatisticsRaw struct {
		json.RawMessage `json:"views"`
	} `json:"pageStatistics"`
	PageStatistics  map[string]PageViews
	Country         string `json:"country"`
	Function        string `json:"function"`
	Industry        string `json:"industry"`
	Region          string `json:"region"`
	Seniority       string `json:"seniority"`
	StaffCountRange string `json:"staffCountRange"`
}

type LifetimeTotalPageStatistics struct {
	ClicksRaw json.RawMessage `json:"clicks"`
	Clicks    map[string]map[string]int64
	ViewsRaw  json.RawMessage `json:"views"`
	Views     map[string]PageViews
}

type PageViews struct {
	PageViews int64 `json:"pageViews"`
}

/*
type TotalPageClicks struct {
	MobileCareersPageClicks MobileCareersPageClicks `json:"mobileCareersPageClicks"`
	CareersPageClicks       CareersPageClicks       `json:"careersPageClicks"`
}*/

/*
type MobileCareersPageClicks struct {
	CareersPageJobsClicks       int `json:"careersPageJobsClicks"`
	CareersPagePromoLinksClicks int `json:"careersPagePromoLinksClicks"`
	CareersPageEmployeesClicks  int `json:"careersPageEmployeesClicks"`
}*/

type PageClicks struct {
	CareersPagePromoLinksClicks  int `json:"careersPagePromoLinksClicks"`
	CareersPageBannerPromoClicks int `json:"careersPageBannerPromoClicks"`
	CareersPageJobsClicks        int `json:"careersPageJobsClicks"`
	CareersPageEmployeesClicks   int `json:"careersPageEmployeesClicks"`
}

/*
type TotalPageViews struct {
	AboutPageViews           PageViews `json:"aboutPageViews"`
	MobileAboutPageViews     PageViews `json:"mobileAboutPageViews"`
	DesktopAboutPageViews    PageViews `json:"desktopAboutPageViews"`
	CareersPageViews         PageViews `json:"careersPageViews"`
	MobileCareersPageViews   PageViews `json:"mobileCareersPageViews"`
	DesktopCareersPageViews  PageViews `json:"desktopCareersPageViews"`
	InsightsPageViews        PageViews `json:"insightsPageViews"`
	MobileInsightsPageViews  PageViews `json:"mobileInsightsPageViews"`
	DesktopInsightsPageViews PageViews `json:"desktopInsightsPageViews"`
	JobsPageViews            PageViews `json:"jobsPageViews"`
	MobileJobsPageViews      PageViews `json:"mobileJobsPageViews"`
	DesktopJobsPageViews     PageViews `json:"desktopJobsPageViews"`
	LifeAtPageViews          PageViews `json:"lifeAtPageViews"`
	MobileLifeAtPageViews    PageViews `json:"mobileLifeAtPageViews"`
	DesktopLifeAtPageViews   PageViews `json:"desktopLifeAtPageViews"`
	OverviewPageViews        PageViews `json:"overviewPageViews"`
	DesktopOverviewPageViews PageViews `json:"desktopOverviewPageViews"`
	AllPageViews             PageViews `json:"allPageViews"`
	AllMobilePageViews       PageViews `json:"allMobilePageViews"`
	AllDesktopPageViews      PageViews `json:"allDesktopPageViews"`
	PeoplePageViews          PageViews `json:"peoplePageViews"`
	MobilePeoplePageViews    PageViews `json:"mobilePeoplePageViews"`
	DesktopPeoplePageViews   PageViews `json:"desktopPeoplePageViews"`
	ProductsPageViews        PageViews `json:"productsPageViews"`
	MobileProductsPageViews  PageViews `json:"mobileProductsPageViews"`
	DesktopProductsPageViews PageViews `json:"desktopProductsPageViews"`
}*/

func (os *OrganizationStats) GetLifetimePageStats(organisationID int) (*[]LifetimePageStats, error) {
	values := url.Values{}
	values.Set("q", "organization")
	values.Set("organization", fmt.Sprintf("urn:li:organization:%v", organisationID))

	urlString := fmt.Sprintf("%s/organizationPageStatistics?%s", os.apiURL, values.Encode())
	fmt.Println(urlString)

	pageStatsResponse := LifetimePageStatsResponse{}

	_, err := os.OAuth2().Get(urlString, &pageStatsResponse)
	if err != nil {
		return nil, err
	}

	for i := range pageStatsResponse.Elements {
		err = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByStaffCountRange)
		if err != nil {
			return nil, err
		}
		err = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByFunction)
		if err != nil {
			return nil, err
		}
		err = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].BySeniority)
		if err != nil {
			return nil, err
		}
		err = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByIndustry)
		if err != nil {
			return nil, err
		}
		err = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByRegion)
		if err != nil {
			return nil, err
		}
		err = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByCountry)
		if err != nil {
			return nil, err
		}

		totalPageViews, err := unmarshalPageViews(pageStatsResponse.Elements[i].Totals.ViewsRaw)
		if err != nil {
			return nil, err
		}
		pageStatsResponse.Elements[i].Totals.Views = *totalPageViews

		totalPageClicks, err := unmarshalPageClicks(pageStatsResponse.Elements[i].Totals.ClicksRaw)
		if err != nil {
			return nil, err
		}
		pageStatsResponse.Elements[i].Totals.Clicks = *totalPageClicks
	}

	return &pageStatsResponse.Elements, nil
}

func unmarshalPageViewsSlice(stats *[]LifetimePageStatisticsByType) error {
	if stats == nil {
		return nil
	}

	for j := range *stats {
		pageViews, err := unmarshalPageViews((*stats)[j].PageStatisticsRaw.RawMessage)
		if err != nil {
			return err
		}
		(*stats)[j].PageStatistics = *pageViews
	}

	return nil
}

func unmarshalPageViews(message json.RawMessage) (*map[string]PageViews, error) {
	pageViews := make(map[string]PageViews)
	err := json.Unmarshal(message, &pageViews)
	if err != nil {
		return nil, err
	}

	return &pageViews, nil
}

func unmarshalPageClicks(message json.RawMessage) (*map[string]map[string]int64, error) {
	pageClicks := make(map[string]map[string]int64)
	err := json.Unmarshal(message, &pageClicks)
	if err != nil {
		return nil, err
	}

	return &pageClicks, nil
}
