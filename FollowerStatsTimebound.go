package linkedin

import (
	"fmt"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

type FollowerStatsTimeboundResponse struct {
	Paging   Paging                   `json:"paging"`
	Elements []FollowerStatsTimebound `json:"elements"`
}

type FollowerStatsTimebound struct {
	TimeRange            TimeRange     `json:"timeRange"`
	FollowerGains        FollowerGains `json:"followerGains"`
	OrganizationalEntity string        `json:"organizationalEntity"`
}

type FollowerGains struct {
	OrganicFollowerGain int64 `json:"organicFollowerGain"`
	PaidFollowerGain    int64 `json:"paidFollowerGain"`
}

func (service *Service) GetFollowerStatsTimebound(organisationID int, startDateUnix int64, endDateUnix int64) (*[]FollowerStatsTimebound, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	followerStatsResponse := FollowerStatsTimeboundResponse{}

	requestConfig := oauth2.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizationalEntityFollowerStatistics?%s", values.Encode())),
		ResponseModel: &followerStatsResponse,
	}
	_, _, e := service.oAuth2.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &followerStatsResponse.Elements, nil
}
