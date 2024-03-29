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

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlV2(fmt.Sprintf("networkSizes/urn:li:organization:%v?edgeType=CompanyFollowedByMember", organizationId)),
		ResponseModel: &organizationNetworkSizes,
	}
	_, _, e := service.versionedHttpRequest(&requestConfig, linkedInVersion)
	if e != nil {
		return nil, e
	}

	return &organizationNetworkSizes, nil
}
