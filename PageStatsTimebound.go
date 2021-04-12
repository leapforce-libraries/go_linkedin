package linkedin

import (
	"fmt"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
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

func (service *Service) GetPageStatsTimebound(organizationID int64, startDateUnix int64, endDateUnix int64) (*[]PageStatsTimebound, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organization")
	values.Set("organization", fmt.Sprintf("urn:li:organization:%v", organizationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	pageStatsResponse := PageStatsTimeboundResponse{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizationPageStatistics?%s", values.Encode())),
		ResponseModel: &pageStatsResponse,
	}
	_, _, e := service.oAuth2.Get(&requestConfig)
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
