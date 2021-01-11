package linkedin

import (
	"fmt"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type PageStatsTimeboundResponse struct {
	Paging   Paging               `json:"paging"`
	Elements []PageStatsTimebound `json:"elements"`
}

type PageStatsTimebound struct {
	TotalPageStatistics TotalPageStatistics `json:"totalPageStatistics"`
	TimeRange           TimeRange           `json:"timeRange"`
	Organization        string              `json:"organization"`
}

func (service *Service) GetPageStatsTimebound(organisationID int, startDateUnix int64, endDateUnix int64) (*[]PageStatsTimebound, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organization")
	values.Set("organization", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	urlString := fmt.Sprintf("%s/organizationPageStatistics?%s", service.BaseURL(), values.Encode())
	//fmt.Println(urlString)

	pageStatsResponse := PageStatsTimeboundResponse{}

	_, _, e := service.OAuth2().Get(urlString, &pageStatsResponse, nil)
	if e != nil {
		return nil, e
	}

	for i := range pageStatsResponse.Elements {
		totalPageViews, e := unmarshalPageViews(pageStatsResponse.Elements[i].TotalPageStatistics.ViewsRaw)
		if e != nil {
			return nil, e
		}
		pageStatsResponse.Elements[i].TotalPageStatistics.Views = *totalPageViews

		totalPageClicks, e := unmarshalPageClicks(pageStatsResponse.Elements[i].TotalPageStatistics.ClicksRaw)
		if e != nil {
			return nil, e
		}
		pageStatsResponse.Elements[i].TotalPageStatistics.Clicks = *totalPageClicks
	}

	return &pageStatsResponse.Elements, nil
}
