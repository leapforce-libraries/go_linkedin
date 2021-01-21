package linkedin

import (
	"fmt"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

type ShareStatsTimeboundResponse struct {
	Paging   Paging                `json:"paging"`
	Elements []ShareStatsTimebound `json:"elements"`
}

type ShareStatsTimebound struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	TimeRange            TimeRange            `json:"timeRange"`
	OrganizationalEntity string               `json:"organizationalEntity"`
	Share                *string              `json:"share"`
}

func (service *Service) GetShareStatsTimebound(organisationID int, startDateUnix int64, endDateUnix int64, shareIDs *[]string) (*[]ShareStatsTimebound, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	if shareIDs != nil {
		for index, shareID := range *shareIDs {
			values.Set(fmt.Sprintf("shares[%v]", index), fmt.Sprintf("urn:li:share:%s", shareID))
		}
	}

	shareStatsResponse := ShareStatsTimeboundResponse{}

	requestConfig := oauth2.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizationalEntityShareStatistics?%s", values.Encode())),
		ResponseModel: &shareStatsResponse,
	}
	_, _, e := service.oAuth2.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &shareStatsResponse.Elements, nil
}
