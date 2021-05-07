package linkedin

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Organization struct {
	VanityName       string `json:"vanityName"`
	LocalizedName    string `json:"localizedName"`
	LocalizedWebsite string `json:"localizedWebsite"`
}

func (service *Service) GetOrganization(organizationID int64) (*Organization, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	organization := Organization{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizations/%v", organizationID)),
		ResponseModel: &organization,
	}
	_, _, e := service.oAuth2Service.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &organization, nil
}
