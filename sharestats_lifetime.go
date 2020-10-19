package linkedin

import (
	"fmt"
	"net/url"
)

type LifetimeShareStatsResponse struct {
	Paging   Paging               `json:"paging"`
	Elements []LifetimeShareStats `json:"elements"`
}

type LifetimeShareStats struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
}

func (li *LinkedIn) GetLifetimeShareStats(organisationID int) (*[]LifetimeShareStats, error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))

	urlString := fmt.Sprintf("%s/organizationalEntityShareStatistics?%s", apiURL, values.Encode())
	//fmt.Println(urlString)

	shareStatsResponse := LifetimeShareStatsResponse{}

	_, err := li.OAuth2().Get(urlString, &shareStatsResponse)
	if err != nil {
		return nil, err
	}

	return &shareStatsResponse.Elements, nil
}
