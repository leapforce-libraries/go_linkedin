package linkedin

import (
	"fmt"

	types "github.com/Leapforce-nl/go_types"
)

type Organization struct {
	VanityName       string `json:"vanityName"`
	LocalizedName    string `json:"localizedName"`
	LocalizedWebsite string `json:"localizedWebsite"`
}

func (li *LinkedIn) GetOrganization(organizationID int) (*Organization, error) {
	if li == nil {
		return nil, &types.ErrorString{"LinkedIn pointer is nil"}
	}

	urlString := fmt.Sprintf("%s/organizations/%v", li.BaseURL(), organizationID)
	//fmt.Println(urlString)

	organization := Organization{}

	_, err := li.OAuth2().Get(urlString, &organization)
	if err != nil {
		return nil, err
	}

	return &organization, nil
}
