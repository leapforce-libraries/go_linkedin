package linkedin

import (
	"fmt"
	"net/http"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	oauth2 "github.com/Leapforce-nl/go_oauth2"
)

const (
	apiName         string = "LinkedIn"
	apiURL          string = "https://api.linkedin.com"
	apiVersion      string = "v2"
	authURL         string = "https://www.linkedin.com/oauth/v2/authorization"
	tokenURL        string = "https://www.linkedin.com/oauth/v2/accessToken"
	tokenHttpMethod string = http.MethodGet
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
	IsLive       bool
}

// NewLinkedIn return new instance of LinkedIn struct
//
func NewLinkedIn(params NewLinkedInParams) (*LinkedIn, error) {
	oa := oauth2.NewOAuth(apiName, params.ClientID, params.ClientSecret, params.Scope, redirectURL, authURL, tokenURL, tokenHttpMethod, params.BigQuery, params.IsLive)
	li := LinkedIn{oa}
	return &li, nil
}

func (li *LinkedIn) OAuth2() *oauth2.OAuth2 {
	return li.oAuth2
}

func (li *LinkedIn) BaseURL() string {
	return fmt.Sprintf("%s/%s", apiURL, apiVersion)
}
