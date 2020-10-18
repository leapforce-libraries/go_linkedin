package linkedin

import (
	oauth2 "github.com/Leapforce-nl/go_oauth2"
)

type OrganizationStats struct {
	apiURL string
	oAuth2 *oauth2.OAuth2
}

func NewOrganizationStats(apiURL string, oa *oauth2.OAuth2) *OrganizationStats {
	return &OrganizationStats{apiURL, oa}
}

func (os *OrganizationStats) OAuth2() *oauth2.OAuth2 {
	return os.oAuth2
}

type Paging struct {
	Count int      `json:"count"`
	Start int      `json:"start"`
	Links []string `json:"links"`
}

type TimeRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}
