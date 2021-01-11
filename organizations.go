package linkedin

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
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

	urlString := fmt.Sprintf("%s/organizations/%v", service.BaseURL(), organizationID)
	//fmt.Println(urlString)

	organization := Organization{}

	_, _, e := service.OAuth2().Get(urlString, &organization, nil)
	if e != nil {
		return nil, e
	}

	return &organization, nil
}
