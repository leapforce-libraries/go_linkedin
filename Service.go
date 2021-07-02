package linkedin

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	google "github.com/leapforce-libraries/go_google"
	bigquery "github.com/leapforce-libraries/go_google/bigquery"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName         string = "LinkedIn"
	apiURL          string = "https://api.linkedin.com/v2"
	authURL         string = "https://www.linkedin.com/oauth/v2/authorization"
	tokenURL        string = "https://www.linkedin.com/oauth/v2/accessToken"
	tokenHTTPMethod string = http.MethodGet
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// LinkedIn stores LinkedIn configuration
//
type Service struct {
	clientID      string
	oAuth2Service *oauth2.Service
}

type ServiceConfig struct {
	ClientID        string
	ClientSecret    string
	BigQueryService *bigquery.Service
}

// NewService return new instance of LinkedIn struct
//
func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	getTokenFunction := func() (*oauth2.Token, *errortools.Error) {
		return google.GetToken(apiName, serviceConfig.ClientID, serviceConfig.BigQueryService)
	}

	saveTokenFunction := func(token *oauth2.Token) *errortools.Error {
		return google.SaveToken(apiName, serviceConfig.ClientID, token, serviceConfig.BigQueryService)
	}

	oAuth2ServiceConfig := oauth2.ServiceConfig{
		ClientID:          serviceConfig.ClientID,
		ClientSecret:      serviceConfig.ClientSecret,
		RedirectURL:       redirectURL,
		AuthURL:           authURL,
		TokenURL:          tokenURL,
		TokenHTTPMethod:   tokenHTTPMethod,
		GetTokenFunction:  &getTokenFunction,
		SaveTokenFunction: &saveTokenFunction,
	}

	oAuth2Service, e := oauth2.NewService(&oAuth2ServiceConfig)
	if e != nil {
		return nil, e
	}

	return &Service{serviceConfig.ClientID, oAuth2Service}, nil
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) InitToken(scope string) *errortools.Error {
	return service.oAuth2Service.InitToken(scope)
}

func (service Service) APIName() string {
	return apiName
}

func (service Service) APIKey() string {
	return service.clientID
}

func (service Service) APICallCount() int64 {
	return service.oAuth2Service.APICallCount()
}

func (service *Service) APIReset() {
	service.oAuth2Service.APIReset()
}
