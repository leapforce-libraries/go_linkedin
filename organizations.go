package linkedin

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type OrganizationsResponse struct {
	Paging   Paging         `json:"paging"`
	Elements []Organization `json:"elements"`
}

type Organization struct {
	VanityName    string `json:"vanityName"`
	LocalizedName string `json:"localizedName"`
	Name          struct {
		Localized struct {
			EnUS string `json:"en_US"`
		} `json:"localized"`
		PreferredLocale struct {
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"preferredLocale"`
	} `json:"name"`
	PrimaryOrganizationType string        `json:"primaryOrganizationType"`
	Locations               []interface{} `json:"locations"`
	Id                      int           `json:"id"`
	LocalizedWebsite        string        `json:"localizedWebsite"`
	LogoV2                  struct {
		Cropped  string `json:"cropped"`
		Original string `json:"original"`
		CropInfo struct {
			X      int `json:"x"`
			Width  int `json:"width"`
			Y      int `json:"y"`
			Height int `json:"height"`
		} `json:"cropInfo"`
	} `json:"logoV2"`
}

func (service *Service) GetOrganization(organizationId int64) (*Organization, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	organization := Organization{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlRest(fmt.Sprintf("organizations/%v", organizationId)),
		ResponseModel: &organization,
	}
	_, _, e := service.versionedHttpRequest(&requestConfig, nil)
	if e != nil {
		return nil, e
	}

	return &organization, nil
}

func (service *Service) FindOrganizationByVanityName(vanityName string) (*[]Organization, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var organizationsResponse OrganizationsResponse

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlRest(fmt.Sprintf("organizations?q=vanityName&vanityName=%s", vanityName)),
		ResponseModel: &organizationsResponse,
	}
	_, _, e := service.versionedHttpRequest(&requestConfig, nil)
	if e != nil {
		return nil, e
	}

	return &organizationsResponse.Elements, nil
}
