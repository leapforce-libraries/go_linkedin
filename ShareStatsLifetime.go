package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type ShareStatsLifetimeResponse struct {
	Paging   Paging               `json:"paging"`
	Elements []ShareStatsLifetime `json:"elements"`
}

type ShareStatsLifetime struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
	Share                *string              `json:"share"`
}

func (service *Service) GetShareStatsLifetime(organisationID int, shareIDs *[]string) (*[]ShareStatsLifetime, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))

	if shareIDs != nil {
		for index, shareID := range *shareIDs {
			values.Set(fmt.Sprintf("shares[%v]", index), fmt.Sprintf("urn:li:share:%s", shareID))
		}
	}

	shareStatsResponse := ShareStatsLifetimeResponse{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizationalEntityShareStatistics?%s", values.Encode())),
		ResponseModel: &shareStatsResponse,
	}
	_, _, e := service.oAuth2.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &shareStatsResponse.Elements, nil
}
