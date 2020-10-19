package linkedin

import (
	"fmt"
	"net/url"
)

type LifetimeFollowerStatsResponse struct {
	Paging   Paging                  `json:"paging"`
	Elements []LifetimeFollowerStats `json:"elements"`
}

type LifetimeFollowerStats struct {
	ByAssociationType    []LifetimeFollowerCountsByType `json:"followerCountsByAssociationType"`
	ByStaffCountRange    []LifetimeFollowerCountsByType `json:"followerCountsByStaffCountRange"`
	ByFunction           []LifetimeFollowerCountsByType `json:"followerCountsByFunction"`
	BySeniority          []LifetimeFollowerCountsByType `json:"followerCountsBySeniority"`
	ByIndustry           []LifetimeFollowerCountsByType `json:"followerCountsByIndustry"`
	ByRegion             []LifetimeFollowerCountsByType `json:"followerCountsByRegion"`
	ByCountry            []LifetimeFollowerCountsByType `json:"followerCountsByCountry"`
	OrganizationalEntity string                         `json:"organizationalEntity"`
}

type LifetimeFollowerCountsByType struct {
	FollowerCounts  FollowerCounts `json:"followerCounts"`
	AssociationType string         `json:"associationType"`
	Country         string         `json:"country"`
	Function        string         `json:"function"`
	Industry        string         `json:"industry"`
	Region          string         `json:"region"`
	Seniority       string         `json:"seniority"`
	StaffCountRange string         `json:"staffCountRange"`
}

type FollowerCounts struct {
	OrganicFollowerCount int64 `json:"organicFollowerCount"`
	PaidFollowerCount    int64 `json:"paidFollowerCount"`
}

func (li *LinkedIn) GetLifetimeFollowerStats(organisationID int) (*[]LifetimeFollowerStats, error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))

	urlString := fmt.Sprintf("%s/organizationalEntityFollowerStatistics?%s", li.BaseURL(), values.Encode())
	//fmt.Println(urlString)

	followerStatsResponse := LifetimeFollowerStatsResponse{}

	_, err := li.OAuth2().Get(urlString, &followerStatsResponse)
	if err != nil {
		return nil, err
	}

	return &followerStatsResponse.Elements, nil
}
