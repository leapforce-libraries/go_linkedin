package linkedin

import (
	"fmt"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
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

	urlString := fmt.Sprintf("%s/organizationalEntityShareStatistics?%s", service.BaseURL(), values.Encode())
	//fmt.Println(urlString)

	shareStatsResponse := ShareStatsTimeboundResponse{}

	_, _, e := service.OAuth2().Get(urlString, &shareStatsResponse, nil)
	if e != nil {
		return nil, e
	}

	return &shareStatsResponse.Elements, nil
}
