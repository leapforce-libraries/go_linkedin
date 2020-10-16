package linkedin

import (
	"net/http"

	bigquerytools "github.com/Leapforce-nl/go_bigquerytools"
	linkedin_os "github.com/Leapforce-nl/go_linkedin/organizationstats"
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
	OrganizationStats *linkedin_os.OrganizationStats
}

type NewLinkedInParams struct {
	ClientID     string
	ClientSecret string
	Scope        string
	BigQuery     *bigquerytools.BigQuery
	IsLive       bool
}

// methods
//
func NewLinkedIn(params NewLinkedInParams) (*LinkedIn, error) {
	li := LinkedIn{}
	oa := oauth2.NewOAuth(apiName, params.ClientID, params.ClientSecret, params.Scope, redirectURL, authURL, tokenURL, tokenHttpMethod, params.BigQuery, params.IsLive)
	li.OrganizationStats = linkedin_os.NewOrganizationStats(apiURL, oa)
	return &li, nil
}
