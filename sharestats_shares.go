package linkedin

import (
	"fmt"
	"net/url"
)

type ShareShareStatsResponse struct {
	Paging   Paging            `json:"paging"`
	Elements []ShareShareStats `json:"elements"`
}

type ShareShareStats struct {
	TotalShareStatistics TotalShareStatistics `json:"totalShareStatistics"`
	OrganizationalEntity string               `json:"organizationalEntity"`
	Share                string               `json:"share"`
}

func (li *LinkedIn) GetShareShareStats(organisationID int, shareIDs []string) (*[]ShareShareStats, error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))

	for index, shareID := range shareIDs {
		values.Set(fmt.Sprintf("shares[%v]", index), fmt.Sprintf("urn:li:share:%s", shareID))
	}

	urlString := fmt.Sprintf("%s/organizationalEntityShareStatistics?%s", apiURL, values.Encode())
	fmt.Println(urlString)

	shareStatsResponse := ShareShareStatsResponse{}

	_, err := li.OAuth2().Get(urlString, &shareStatsResponse)
	if err != nil {
		return nil, err
	}

	return &shareStatsResponse.Elements, nil
}
