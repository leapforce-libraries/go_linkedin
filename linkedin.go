package LinkedIn

import (
	"net/http"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	oauth2 "github.com/Leapforce-nl/go_oauth2"
)

const (
	apiName         string = "LinkedIn"
	apiURL          string = "https://api.linkedin.com/v2"
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

// methods
//
func NewLinkedIn(clientID string, clientSecret string, scope string, bigQuery *bigquerytools.BigQuery, isLive bool) (*LinkedIn, error) {
	li := LinkedIn{}
	li.oAuth2 = oauth2.NewOAuth(apiName, clientID, clientSecret, scope, redirectURL, authURL, tokenURL, tokenHttpMethod, bigQuery, isLive)
	return &li, nil
}

func (li *LinkedIn) ValidateToken() error {
	return li.oAuth2.ValidateToken()
}

func (li *LinkedIn) InitToken() error {
	return li.oAuth2.InitToken()
}
