package linkedin

import (
	"fmt"
	"net/http"
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
	ByGeoCountry         []LifetimeFollowerCountsByType `json:"followerCountsByGeoCountry"`
	ByFunction           []LifetimeFollowerCountsByType `json:"followerCountsByFunction"`
	ByIndustry           []LifetimeFollowerCountsByType `json:"followerCountsByIndustry"`
	ByGeo                []LifetimeFollowerCountsByType `json:"followerCountsByGeo"`
	BySeniority          []LifetimeFollowerCountsByType `json:"followerCountsBySeniority"`
	ByStaffCountRange    []LifetimeFollowerCountsByType `json:"followerCountsByStaffCountRange"`
	OrganizationalEntity string                         `json:"organizationalEntity"`
}

type LifetimeFollowerCountsByType struct {
	FollowerCounts  FollowerCounts `json:"followerCounts"`
	AssociationType string         `json:"associationType"`
	Geo             string         `json:"geo"`
	Function        string         `json:"function"`
	Industry        string         `json:"industry"`
	Seniority       string         `json:"seniority"`
	StaffCountRange string         `json:"staffCountRange"`
}

type FollowerCounts struct {
	OrganicFollowerCount int64 `json:"organicFollowerCount"`
	PaidFollowerCount    int64 `json:"paidFollowerCount"`
}

func (service *Service) GetFollowerStatsLifetime(organizationId int64) (*[]FollowerStatsLifetime, *errortools.Error) {
	values := url.Values{}
	values.Set("q", "organizationalEntity")
	values.Set("organizationalEntity", fmt.Sprintf("urn:li:organization:%v", organizationId))

	followerStatsResponse := FollowerStatsLifetimeResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlRest(fmt.Sprintf("organizationalEntityFollowerStatistics?%s", values.Encode())),
		ResponseModel: &followerStatsResponse,
	}
	_, _, e := service.versionedHttpRequest(&requestConfig, nil)
	if e != nil {
		return nil, e
	}

	return &followerStatsResponse.Elements, nil
}
