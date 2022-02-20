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
	apiUrl                 string = "https://api.linkedin.com/v2"
	authUrl                string = "https://www.linkedin.com/oauth/v2/authorization"
	tokenUrl               string = "https://www.linkedin.com/oauth/v2/accessToken"
	tokenHttpMethod        string = http.MethodGet
	redirectUrl            string = "http://localhost:8080/oauth/redirect"
	CampaignUrnPrefix      string = "urn:li:sponsoredCampaign:"
	CreativeUrnPrefix      string = "urn:li:sponsoredCreative:"
	InMailContentUrnPrefix string = "urn:li:adInMailContent:"
	countDefault           uint   = 10
	maxUrnsPerCall         uint   = 50
)

// LinkedIn stores LinkedIn configuration
//
type Service struct {
	clientId      string
	oAuth2Service *oauth2.Service
}

type ServiceConfig struct {
	ClientId     string
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
		ClientId:        serviceConfig.ClientId,
		ClientSecret:    serviceConfig.ClientSecret,
		RedirectUrl:     redirectUrl,
		AuthUrl:         authUrl,
		TokenUrl:        tokenUrl,
		TokenHttpMethod: tokenHttpMethod,
		TokenSource:     tokenMap,
	}

	oAuth2Service, e := oauth2.NewService(&oAuth2ServiceConfig)
	if e != nil {
		return nil, e
	}

	return &Service{serviceConfig.ClientId, oAuth2Service}, nil
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service *Service) FromUrn(prefix string, urn string) int64 {
	id, err := strconv.ParseInt(strings.TrimPrefix(urn, prefix), 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func (service *Service) InitToken(scope string, accessType *string, prompt *string, state *string) *errortools.Error {
	return service.oAuth2Service.InitToken(scope, accessType, prompt, state)
}

func (service Service) ApiName() string {
	return apiName
}

func (service Service) ApiKey() string {
	return service.clientId
}

func (service Service) ApiCallCount() int64 {
	return service.oAuth2Service.ApiCallCount()
}

func (service *Service) ApiReset() {
	service.oAuth2Service.ApiReset()
}
