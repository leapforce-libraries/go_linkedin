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

func (li *LinkedIn) GetOrganization(organizationID int) (*Organization, *errortools.Error) {
	if li == nil {
		return nil, errortools.ErrorMessage("LinkedIn pointer is nil")
	}

	urlString := fmt.Sprintf("%s/organizations/%v", li.BaseURL(), organizationID)
	//fmt.Println(urlString)

	organization := Organization{}

	_, _, e := li.OAuth2().Get(urlString, &organization, nil)
	if e != nil {
		return nil, e
	}

	return &organization, nil
}
