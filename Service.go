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
	oAuth2 *oauth2.OAuth2
}

type ServiceConfig struct {
	ClientID     string
	ClientSecret string
	Scope        string
	BigQuery     *bigquery.Service
}

// NewLinkedIn return new instance of LinkedIn struct
//
func NewService(config ServiceConfig) *Service {
	getTokenFunction := func() (*oauth2.Token, *errortools.Error) {
		return google.GetToken(apiName, config.ClientID, config.BigQuery)
	}

	saveTokenFunction := func(token *oauth2.Token) *errortools.Error {
		return google.SaveToken(apiName, config.ClientID, token, config.BigQuery)
	}

	oauth2Config := oauth2.OAuth2Config{
		ClientID:          config.ClientID,
		ClientSecret:      config.ClientSecret,
		RedirectURL:       redirectURL,
		AuthURL:           authURL,
		TokenURL:          tokenURL,
		TokenHTTPMethod:   tokenHTTPMethod,
		GetTokenFunction:  &getTokenFunction,
		SaveTokenFunction: &saveTokenFunction,
	}

	return &Service{oauth2.NewOAuth(oauth2Config)}
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) InitToken() *errortools.Error {
	return service.oAuth2.InitToken()
}
