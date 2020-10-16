package linkedin

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
)

type TimeboundFollowerStatsResponse struct {
	Paging   Paging                   `json:"paging"`
	Elements []TimeboundFollowerStats `json:"elements"`
}

type TimeboundFollowerStats struct {
	TimeRange            TimeRange     `json:"timeRange"`
	FollowerGains        FollowerGains `json:"followerGains"`
	OrganizationalEntity string        `json:"organizationalEntity"`
}

type FollowerGains struct {
	OrganicFollowerGain int `json:"organicFollowerGain"`
	PaidFollowerGain    int `json:"paidFollowerGain"`
}

type TimeRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

func (os *OrganizationStats) GetTimeboundFollowerStats(organisationID int, startDate civil.Date, endDate civil.Date) (*[]TimeboundFollowerStats, error) {
	location, _ := time.LoadLocation("GMT")
	unixStart := startDate.In(location).Unix() * 1000
	unixEnd := endDate.In(location).Unix() * 1000

	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(unixStart, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(unixEnd, 10))

	urlString := fmt.Sprintf("%s/organizationalEntityFollowerStatistics?%s", os.apiURL, values.Encode())
	fmt.Println(urlString)

	followerStatsResponse := TimeboundFollowerStatsResponse{}

	_, err := os.OAuth2().Get(urlString, &followerStatsResponse)
	if err != nil {
		return nil, err
	}

	return &followerStatsResponse.Elements, nil
}
