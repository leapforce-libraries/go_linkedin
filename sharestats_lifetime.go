package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type LifetimeShareStatsResponse struct {
	Paging   Paging               `json:"paging"`
	Elements []LifetimeShareStats `json:"elements"`
}

type LifetimeShareStats struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
}

func (li *LinkedIn) GetLifetimeShareStats(organisationID int) (*[]LifetimeShareStats, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))

	urlString := fmt.Sprintf("%s/organizationalEntityShareStatistics?%s", li.BaseURL(), values.Encode())
	//fmt.Println(urlString)

	shareStatsResponse := LifetimeShareStatsResponse{}

	_, _, e := li.OAuth2().Get(urlString, &shareStatsResponse, nil)
	if e != nil {
		return nil, e
	}

	return &shareStatsResponse.Elements, nil
}
