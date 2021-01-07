package linkedin

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	google "github.com/leapforce-libraries/go_google"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

const (
	APIName         string = "LinkedIn"
	APIURL          string = "https://api.linkedin.com"
	APIVersion      string = "v2"
	AuthURL         string = "https://www.linkedin.com/oauth/v2/authorization"
	TokenURL        string = "https://www.linkedin.com/oauth/v2/accessToken"
	TokenHTTPMethod string = http.MethodGet
	RedirectURL     string = "http://localhost:8080/oauth/redirect"
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
	BigQuery     *google.BigQuery
}

// NewLinkedIn return new instance of LinkedIn struct
//
func NewLinkedIn(params NewLinkedInParams) *LinkedIn {
	getTokenFunction := func() (*oauth2.Token, *errortools.Error) {
		return google.GetToken(APIName, params.ClientID, params.BigQuery)
	}

	saveTokenFunction := func(token *oauth2.Token) *errortools.Error {
		return google.SaveToken(APIName, params.ClientID, token, params.BigQuery)
	}

	config := oauth2.OAuth2Config{
		ClientID:          params.ClientID,
		ClientSecret:      params.ClientSecret,
		RedirectURL:       RedirectURL,
		AuthURL:           AuthURL,
		TokenURL:          TokenURL,
		TokenHTTPMethod:   TokenHTTPMethod,
		GetTokenFunction:  &getTokenFunction,
		SaveTokenFunction: &saveTokenFunction,
	}

	return &LinkedIn{oauth2.NewOAuth(config)}
}

func (li *LinkedIn) OAuth2() *oauth2.OAuth2 {
	return li.oAuth2
}

func (li *LinkedIn) BaseURL() string {
	return fmt.Sprintf("%s/%s", APIURL, APIVersion)
}

func (li *LinkedIn) InitToken() *errortools.Error {
	return li.oAuth2.InitToken()
}
