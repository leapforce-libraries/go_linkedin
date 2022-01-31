package linkedin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	gcs "github.com/leapforce-libraries/go_googlecloudstorage"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	go_tokenmap "github.com/leapforce-libraries/go_oauth2/tokenmap"
)

const (
	apiName                string = "LinkedIn"
	apiURL                 string = "https://api.linkedin.com/v2"
	authURL                string = "https://www.linkedin.com/oauth/v2/authorization"
	tokenURL               string = "https://www.linkedin.com/oauth/v2/accessToken"
	tokenHTTPMethod        string = http.MethodGet
	redirectURL            string = "http://localhost:8080/oauth/redirect"
	CampaignURNPrefix      string = "urn:li:sponsoredCampaign:"
	CreativeURNPrefix      string = "urn:li:sponsoredCreative:"
	InMailContentURNPrefix string = "urn:li:adInMailContent:"
	countDefault           uint   = 10
	maxURNsPerCall         uint   = 50
)

// LinkedIn stores LinkedIn configuration
//
type Service struct {
	clientID      string
	oAuth2Service *oauth2.Service
}

type ServiceConfig struct {
	ClientID     string
	ClientSecret string
	CredMap      *gcs.Map
}

// NewService return new instance of LinkedIn struct
//
func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}
	if serviceConfig.CredMap == nil {
		return nil, errortools.ErrorMessage("CredMap must not be a nil pointer")
	}

	tokenMap, e := go_tokenmap.NewTokenMap(serviceConfig.CredMap)
	if e != nil {
		return nil, e
	}

	oAuth2ServiceConfig := oauth2.ServiceConfig{
		ClientID:        serviceConfig.ClientID,
		ClientSecret:    serviceConfig.ClientSecret,
		RedirectURL:     redirectURL,
		AuthURL:         authURL,
		TokenURL:        tokenURL,
		TokenHTTPMethod: tokenHTTPMethod,
		TokenSource:     tokenMap,
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

func (service *Service) FromURN(prefix string, urn string) int64 {
	id, err := strconv.ParseInt(strings.TrimPrefix(urn, prefix), 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func (service *Service) InitToken(scope string, accessType *string, prompt *string, state *string) *errortools.Error {
	return service.oAuth2Service.InitToken(scope, accessType, prompt, state)
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
