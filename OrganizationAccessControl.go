package linkedin

import (
	"fmt"
	"net/http"
	"net/url"

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

	var count = 100
	var start = 0

	var organizationAcls []OrganizationAcl

	var values url.Values
	values.Set("q", "roleAssignee")
	values.Set("count", fmt.Sprintf("%v", count))

	for {
		values.Set("start", fmt.Sprintf("%v", start))

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

		organizationAcls = append(organizationAcls, response.Elements...)

		if len(response.Elements) < count {
			break
		}

		start += count
	}

	return &organizationAcls, nil
}
