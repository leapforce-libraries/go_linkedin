package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type UGCPostStatsLifetimeResponse struct {
	Paging   Paging                 `json:"paging"`
	Elements []UGCPostStatsLifetime `json:"elements"`
}

type UGCPostStatsLifetime struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
	UGCPost              *string              `json:"ugcPost"`
}

func (service *Service) GetUGCPostStatsLifetime(organizationID int64, ugcPostIDs *[]string) (*[]UGCPostStatsLifetime, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organizationID))

	if ugcPostIDs != nil {
		for index, ugcPostID := range *ugcPostIDs {
			values.Set(fmt.Sprintf("ugcPosts[%v]", index), ugcPostID)
		}
	}

	ugcPostStatsResponse := UGCPostStatsLifetimeResponse{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizationalEntityShareStatistics?%s", values.Encode())),
		ResponseModel: &ugcPostStatsResponse,
	}
	_, _, e := service.oAuth2Service.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &ugcPostStatsResponse.Elements, nil
}
