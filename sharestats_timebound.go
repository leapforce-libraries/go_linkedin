package linkedin

import (
	"fmt"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type TimeboundShareStatsResponse struct {
	Paging   Paging                `json:"paging"`
	Elements []TimeboundShareStats `json:"elements"`
}

type TimeboundShareStats struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	TimeRange            TimeRange            `json:"timeRange"`
	OrganizationalEntity string               `json:"organizationalEntity"`
}

func (li *LinkedIn) GetTimeboundShareStats(organisationID int, startDateUnix int64, endDateUnix int64) (*[]TimeboundShareStats, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))
	values.Set("timeIntervals.timeGranularityType", "DAY")
	values.Set("timeIntervals.timeRange.start", strconv.FormatInt(startDateUnix, 10))
	values.Set("timeIntervals.timeRange.end", strconv.FormatInt(endDateUnix, 10))

	urlString := fmt.Sprintf("%s/organizationalEntityShareStatistics?%s", li.BaseURL(), values.Encode())
	//fmt.Println(urlString)

	shareStatsResponse := TimeboundShareStatsResponse{}

	_, _, e := li.OAuth2().Get(urlString, &shareStatsResponse, nil)
	if e != nil {
		return nil, e
	}

	return &shareStatsResponse.Elements, nil
}
