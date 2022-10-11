package linkedin

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type OrganizationNetworkSizes struct {
	FirstDegreeSize int64 `json:"firstDegreeSize"`
}

func (service *Service) GetOrganizationNetworkSizes(organizationId int64, linkedInVersion *string) (*OrganizationNetworkSizes, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	organizationNetworkSizes := OrganizationNetworkSizes{}

	headers := http.Header{}
	version := defaultLinkedInVersion
	if linkedInVersion != nil {
		version = *linkedInVersion
	}
	headers.Set(linkedInVersionHeader, version)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodGet,
		Url:               service.urlRest(fmt.Sprintf("networkSizes/urn:li:organization:%v?edgeType=CompanyFollowedByMember", organizationId)),
		ResponseModel:     &organizationNetworkSizes,
		NonDefaultHeaders: &headers,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &organizationNetworkSizes, nil
}
