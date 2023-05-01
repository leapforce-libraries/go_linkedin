package linkedin

import (
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type OrganizationAclResponse struct {
	Paging   Paging            `json:"paging"`
	Elements []OrganizationAcl `json:"elements"`
}

type OrganizationAcl struct {
	RoleAssignee string `json:"roleAssignee"`
	State        string `json:"state"`
	Role         string `json:"role"`
	Organization string `json:"organization"`
}

func (service *Service) GetOrganizationAcls() (*[]OrganizationAcl, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	response := OrganizationAclResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlRest("organizationAcls?q=roleAssignee"),
		ResponseModel: &response,
	}
	_, _, e := service.versionedHttpRequest(&requestConfig, nil)
	if e != nil {
		return nil, e
	}

	return &response.Elements, nil
}
