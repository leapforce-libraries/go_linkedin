package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type ShareShareStatsResponse struct {
	Paging   Paging            `json:"paging"`
	Elements []ShareShareStats `json:"elements"`
}

type ShareShareStats struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
	Share                string               `json:"share"`
}

func (service *Service) GetShareShareStats(organisationID int, shareIDs []string) (*[]ShareShareStats, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:service:organization:%v", organisationID))

	for index, shareID := range shareIDs {
		values.Set(fmt.Sprintf("shares[%v]", index), fmt.Sprintf("urn:service:share:%s", shareID))
	}

	urlString := fmt.Sprintf("%s/organizationalEntityShareStatistics?%s", service.BaseURL(), values.Encode())
	//fmt.Println(urlString)

	shareStatsResponse := ShareShareStatsResponse{}

	_, _, e := service.OAuth2().Get(urlString, &shareStatsResponse, nil)
	if e != nil {
		return nil, e
	}

	return &shareStatsResponse.Elements, nil
}
