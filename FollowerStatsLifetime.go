package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type FollowerStatsLifetimeResponse struct {
	Paging   Paging                  `json:"paging"`
	Elements []FollowerStatsLifetime `json:"elements"`
}

type FollowerStatsLifetime struct {
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

func (service *Service) GetFollowerStatsLifetime(organisationID int) (*[]FollowerStatsLifetime, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organisationID))

	followerStatsResponse := FollowerStatsLifetimeResponse{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("organizationalEntityFollowerStatistics?%s", values.Encode())),
		ResponseModel: &followerStatsResponse,
	}
	_, _, e := service.oAuth2.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &followerStatsResponse.Elements, nil
}
