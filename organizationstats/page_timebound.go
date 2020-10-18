package linkedin

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
)

type TimeboundPageStatsResponse struct {
	Paging   Paging               `json:"paging"`
	Elements []TimeboundPageStats `json:"elements"`
}

type TimeboundPageStats struct {
	TotalPageStatistics TotalPageStatistics `json:"totalPageStatistics"`
	TimeRange           TimeRange           `json:"timeRange"`
	Organization        string              `json:"organization"`
}

func (os *OrganizationStats) GetTimeboundPageStats(organisationID int, startDate civil.Date, endDate civil.Date) (*[]TimeboundPageStats, error) {
	location, _ := time.LoadLocation("GMT")
	unixStart := startDate.In(location).Unix() * 1000
	unixEnd := endDate.In(location).Unix() * 1000

	values := url.Values{}
	values.Set("q", "organization")
	values.Set("organization", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(unixStart, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(unixEnd, 10))

	urlString := fmt.Sprintf("%s/organizationPageStatistics?%s", os.apiURL, values.Encode())
	fmt.Println(urlString)

	pageStatsResponse := TimeboundPageStatsResponse{}

	_, err := os.OAuth2().Get(urlString, &pageStatsResponse)
	if err != nil {
		return nil, err
	}

	for i := range pageStatsResponse.Elements {
		totalPageViews, err := unmarshalPageViews(pageStatsResponse.Elements[i].TotalPageStatistics.ViewsRaw)
		if err != nil {
			return nil, err
		}
		pageStatsResponse.Elements[i].TotalPageStatistics.Views = *totalPageViews

		totalPageClicks, err := unmarshalPageClicks(pageStatsResponse.Elements[i].TotalPageStatistics.ClicksRaw)
		if err != nil {
			return nil, err
		}
		pageStatsResponse.Elements[i].TotalPageStatistics.Clicks = *totalPageClicks
	}

	return &pageStatsResponse.Elements, nil
}
