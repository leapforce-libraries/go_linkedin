package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

type PageStatsLifetimeResponse struct {
	Paging   Paging              `json:"paging"`
	Elements []PageStatsLifetime `json:"elements"`
}

type PageStatsLifetime struct {
	ByStaffCountRange []LifetimePageStatisticsByType `json:"pageStatisticsByStaffCountRange"`
	ByFunction        []LifetimePageStatisticsByType `json:"pageStatisticsByFunction"`
	BySeniority       []LifetimePageStatisticsByType `json:"pageStatisticsBySeniority"`
	ByIndustry        []LifetimePageStatisticsByType `json:"pageStatisticsByIndustry"`
	ByRegion          []LifetimePageStatisticsByType `json:"pageStatisticsByRegion"`
	ByCountry         []LifetimePageStatisticsByType `json:"pageStatisticsByCountry"`
	Totals            TotalPageStatistics            `json:"totalPageStatistics"`
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

func (service *Service) GetPageStatsLifetime(organisationID int) (*[]PageStatsLifetime, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organization")
	values.Set("organization", fmt.Sprintf("urn:li:organization:%v", organisationID))

	pageStatsResponse := PageStatsLifetimeResponse{}

	requestConfig := oauth2.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizationPageStatistics?%s", values.Encode())),
		ResponseModel: &pageStatsResponse,
	}
	_, _, e := service.oAuth2.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	for i := range pageStatsResponse.Elements {
		e = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByStaffCountRange)
		if e != nil {
			return nil, e
		}
		e = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByFunction)
		if e != nil {
			return nil, e
		}
		e = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].BySeniority)
		if e != nil {
			return nil, e
		}
		e = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByIndustry)
		if e != nil {
			return nil, e
		}
		e = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByRegion)
		if e != nil {
			return nil, e
		}
		e = unmarshalPageViewsSlice(&pageStatsResponse.Elements[i].ByCountry)
		if e != nil {
			return nil, e
		}

		totalPageViews, e := unmarshalPageViews(pageStatsResponse.Elements[i].Totals.ViewsRaw)
		if e != nil {
			return nil, e
		}
		pageStatsResponse.Elements[i].Totals.Views = *totalPageViews

		totalPageClicks, e := unmarshalPageClicks(pageStatsResponse.Elements[i].Totals.ClicksRaw)
		if e != nil {
			return nil, e
		}
		pageStatsResponse.Elements[i].Totals.Clicks = *totalPageClicks
	}

	return &pageStatsResponse.Elements, nil
}

func unmarshalPageViewsSlice(stats *[]LifetimePageStatisticsByType) *errortools.Error {
	if stats == nil {
		return nil
	}

	for j := range *stats {
		pageViews, e := unmarshalPageViews((*stats)[j].PageStatisticsRaw.RawMessage)
		if e != nil {
			return e
		}
		(*stats)[j].PageStatistics = *pageViews
	}

	return nil
}

func unmarshalPageViews(message json.RawMessage) (*map[string]PageViews, *errortools.Error) {
	pageViews := make(map[string]PageViews)
	err := json.Unmarshal(message, &pageViews)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	return &pageViews, nil
}

func unmarshalPageClicks(message json.RawMessage) (*map[string]map[string]int64, *errortools.Error) {
	pageClicks_ := make(map[string]json.RawMessage)
	err := json.Unmarshal(message, &pageClicks_)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	pageClicks := make(map[string]map[string]int64)

	for key, value := range pageClicks_ {
		pc := make(map[string]int64)
		err := json.Unmarshal(value, &pc)
		if err == nil {
			pageClicks[key] = pc
		}
	}

	return &pageClicks, nil
}
