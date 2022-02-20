package linkedin

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Organization struct {
	VanityName       string `json:"vanityName"`
	LocalizedName    string `json:"localizedName"`
	LocalizedWebsite string `json:"localizedWebsite"`
}

func (service *Service) GetOrganization(organizationId int64) (*Organization, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	organization := Organization{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("organizations/%v", organizationId)),
		ResponseModel: &organization,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &organization, nil
}
