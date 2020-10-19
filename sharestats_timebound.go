package linkedin

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	general "github.com/Leapforce-nl/go_linkedin/general"
)

type TimeboundShareStatsResponse struct {
	Paging   general.Paging        `json:"paging"`
	Elements []TimeboundShareStats `json:"elements"`
}

type TimeboundShareStats struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	TimeRange            general.TimeRange    `json:"timeRange"`
	OrganizationalEntity string               `json:"organizationalEntity"`
}

func (li *LinkedIn) GetTimeboundShareStats(organisationID int, startDate civil.Date, endDate civil.Date) (*[]TimeboundShareStats, error) {
	location, _ := time.LoadLocation("GMT")
	unixStart := startDate.In(location).Unix() * 1000
	unixEnd := endDate.In(location).Unix() * 1000

	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(unixStart, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(unixEnd, 10))

	urlString := fmt.Sprintf("%s/organizationalEntityShareStatistics?%s", apiURL, values.Encode())
	//fmt.Println(urlString)

	shareStatsResponse := TimeboundShareStatsResponse{}

	_, err := li.OAuth2().Get(urlString, &shareStatsResponse)
	if err != nil {
		return nil, err
	}

	return &shareStatsResponse.Elements, nil
}
