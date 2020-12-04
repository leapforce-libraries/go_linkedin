package linkedin

import (
	"fmt"
	"net/http"

	bigquerytools "github.com/leapforce-libraries/go_bigquerytools"
	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	apiName         string = "LinkedIn"
	apiURL          string = "https://api.linkedin.com"
	apiVersion      string = "v2"
	authURL         string = "https://www.linkedin.com/oauth/v2/authorization"
	tokenURL        string = "https://www.linkedin.com/oauth/v2/accessToken"
	tokenHTTPMethod string = http.MethodGet
	redirectURL     string = "http://localhost:8080/oauth/redirect"
)

// LinkedIn stores LinkedIn configuration
//
type LinkedIn struct {
	oAuth2 *oauth2.OAuth2
}

type NewLinkedInParams struct {
	ClientID     string
	ClientSecret string
	Scope        string
	BigQuery     *bigquerytools.BigQuery
}

// NewLinkedIn return new instance of LinkedIn struct
//
func NewLinkedIn(params NewLinkedInParams) *LinkedIn {
	config := oauth2.OAuth2Config{
		ApiName:         apiName,
		ClientID:        params.ClientID,
		ClientSecret:    params.ClientSecret,
		Scope:           params.Scope,
		RedirectURL:     redirectURL,
		AuthURL:         authURL,
		TokenURL:        tokenURL,
		TokenHTTPMethod: tokenHTTPMethod,
	}
	oa := oauth2.NewOAuth(config, params.BigQuery)
	li := LinkedIn{oa}
	return &li
}

func (li *LinkedIn) OAuth2() *oauth2.OAuth2 {
	return li.oAuth2
}

func (li *LinkedIn) BaseURL() string {
	return fmt.Sprintf("%s/%s", apiURL, apiVersion)
}

func (li *LinkedIn) InitToken() *errortools.Error {
	return li.oAuth2.InitToken()
}
