package linkedin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	go_token "github.com/leapforce-libraries/go_oauth2/token"
	tokensource "github.com/leapforce-libraries/go_oauth2/tokensource"
)

const (
	apiName                string = "LinkedIn"
	apiUrl                 string = "https://api.linkedin.com/v2"
	authUrl                string = "https://www.linkedin.com/oauth/v2/authorization"
	tokenUrl               string = "https://www.linkedin.com/oauth/v2/accessToken"
	tokenHttpMethod        string = http.MethodPost
	defaultRedirectUrl     string = "http://localhost:8080/oauth/redirect"
	CampaignUrnPrefix      string = "urn:li:sponsoredCampaign:"
	CreativeUrnPrefix      string = "urn:li:sponsoredCreative:"
	InMailContentUrnPrefix string = "urn:li:adInMailContent:"
	OrganizationUrnPrefix  string = "urn:li:organization:"
	ShareUrnPrefix         string = "urn:li:share:"
	UgcPostUrnPrefix       string = "urn:li:ugcPost:"
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
	TokenSource  tokensource.TokenSource
	RedirectUrl  *string
}

// NewService return new instance of LinkedIn struct
//
func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	redirectUrl := defaultRedirectUrl
	if serviceConfig.RedirectUrl != nil {
		redirectUrl = *serviceConfig.RedirectUrl
	}

	oAuth2ServiceConfig := oauth2.ServiceConfig{
		ClientId:        serviceConfig.ClientId,
		ClientSecret:    serviceConfig.ClientSecret,
		RedirectUrl:     redirectUrl,
		AuthUrl:         authUrl,
		TokenUrl:        tokenUrl,
		TokenHttpMethod: tokenHttpMethod,
		TokenSource:     serviceConfig.TokenSource,
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

func (service *Service) AuthorizeUrl(scope string, accessType *string, prompt *string, state *string) string {
	return service.oAuth2Service.AuthorizeUrl(scope, accessType, prompt, state)
}

func (service *Service) ValidateToken() (*go_token.Token, *errortools.Error) {
	return service.oAuth2Service.ValidateToken()
}

func (service *Service) GetTokenFromCode(r *http.Request, checkState *func(state string) *errortools.Error) *errortools.Error {
	return service.oAuth2Service.GetTokenFromCode(r, checkState)
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
