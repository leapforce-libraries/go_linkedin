package linkedin

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type IntrospectTokenResponse struct {
	Active       bool   `json:"active"`
	ClientId     string `json:"client_id"`
	AuthorizedAt int    `json:"authorized_at"`
	CreatedAt    int    `json:"created_at"`
	Status       string `json:"status"`
	ExpiresAt    int    `json:"expires_at"`
	Scope        string `json:"scope"`
	AuthType     string `json:"auth_type"`
}

func (service *Service) IntrospectToken(token string) (*IntrospectTokenResponse, *errortools.Error) {
	var response IntrospectTokenResponse

	t := true
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	requestConfig := go_http.RequestConfig{
		Method: http.MethodPost,
		Url:    service.urlOAuth("introspectToken"),
		BodyModel: struct {
			ClientId     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			Token        string `json:"token"`
		}{
			service.clientId,
			service.clientSecret,
			token,
		},
		ResponseModel:      &response,
		XWwwFormUrlEncoded: &t,
		NonDefaultHeaders:  &header,
	}

	httpService, e := go_http.NewService(nil)
	if e != nil {
		return nil, e
	}
	_, _, e = httpService.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response, nil
}
