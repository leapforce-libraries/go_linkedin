package linkedin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
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

func (service *Service) GetShareStatsTimebound(organizationId int64, startDateUnix int64, endDateUnix int64, shareIds *[]string) (*[]ShareStatsTimebound, *http.Response, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organizationId))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	if shareIds != nil {
		for index, shareId := range *shareIds {
			values.Set(fmt.Sprintf("shares[%v]", index), shareId)
		}
	}

	shareStatsResponse := ShareStatsTimeboundResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("organizationalEntityShareStatistics?%s", values.Encode())),
		ResponseModel: &shareStatsResponse,
	}
	_, response, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, response, e
	}

	return &shareStatsResponse.Elements, response, nil
}
