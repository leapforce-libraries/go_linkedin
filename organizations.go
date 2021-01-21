package linkedin

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

type Organization struct {
	VanityName       string `json:"vanityName"`
	LocalizedName    string `json:"localizedName"`
	LocalizedWebsite string `json:"localizedWebsite"`
}

func (service *Service) GetOrganization(organizationID int) (*Organization, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	organization := Organization{}

	requestConfig := oauth2.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizations/%v", organizationID)),
		ResponseModel: &organization,
	}
	_, _, e := service.oAuth2.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &organization, nil
}
