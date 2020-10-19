package linkedin

import (
	"fmt"
	"net/url"
	"strconv"
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
	OrganicFollowerGain int64 `json:"organicFollowerGain"`
	PaidFollowerGain    int64 `json:"paidFollowerGain"`
}

func (li *LinkedIn) GetTimeboundFollowerStats(organisationID int, startDateUnix int64, endDateUnix int64) (*[]TimeboundFollowerStats, error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	urlString := fmt.Sprintf("%s/organizationalEntityFollowerStatistics?%s", li.BaseURL(), values.Encode())
	//fmt.Println(urlString)

	followerStatsResponse := TimeboundFollowerStatsResponse{}

	_, err := li.OAuth2().Get(urlString, &followerStatsResponse)
	if err != nil {
		return nil, err
	}

	return &followerStatsResponse.Elements, nil
}
