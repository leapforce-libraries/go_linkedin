package linkedin

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type UgcPostStatsLifetimeResponse struct {
	Paging   Paging                 `json:"paging"`
	Elements []UgcPostStatsLifetime `json:"elements"`
}

type UgcPostStatsLifetime struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
	UgcPost              *string              `json:"ugcPost"`
}

func (service *Service) GetUgcPostStatsLifetime(organizationId int64, ugcPostIds *[]string) (*[]UgcPostStatsLifetime, *http.Response, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organizationId))

	if ugcPostIds != nil {
		for index, ugcPostId := range *ugcPostIds {
			values.Set(fmt.Sprintf("ugcPosts[%v]", index), ugcPostId)
		}
	}

	ugcPostStatsResponse := UgcPostStatsLifetimeResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("organizationalEntityShareStatistics?%s", values.Encode())),
		ResponseModel: &ugcPostStatsResponse,
	}
	_, response, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, response, e
	}

	return &ugcPostStatsResponse.Elements, response, nil
}
