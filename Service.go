package linkedin

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	go_token "github.com/leapforce-libraries/go_oauth2/token"
	tokensource "github.com/leapforce-libraries/go_oauth2/tokensource"
	"net/http"
	"strconv"
	"strings"
)

const (
	apiName                      string = "LinkedIn"
	apiUrlRest                   string = "https://api.linkedin.com/rest"
	apiUrl                       string = "https://api.linkedin.com"
	authUrl                      string = "https://www.linkedin.com/oauth/v2/authorization"
	tokenUrl                     string = "https://www.linkedin.com/oauth/v2/accessToken"
	linkedInVersionHeader        string = "LinkedIn-Version"
	defaultLinkedInVersion       string = "202304"
	restliProtocolVersionHeader  string = "X-Restli-Protocol-Version"
	defaultRestliProtocolVersion string = "2.0.0"
	tokenHttpMethod              string = http.MethodPost
	defaultRedirectUrl           string = "http://localhost:8080/oauth/redirect"
	AccountUrnPrefix             string = "urn:li:sponsoredAccount:"
	CampaignUrnPrefix            string = "urn:li:sponsoredCampaign:"
	CreativeUrnPrefix            string = "urn:li:sponsoredCreative:"
	InMailContentUrnPrefix       string = "urn:li:adInMailContent:"
	OrganizationUrnPrefix        string = "urn:li:organization:"
	ShareUrnPrefix               string = "urn:li:share:"
	UgcPostUrnPrefix             string = "urn:li:ugcPost:"
	PostUrnPrefix                string = "urn:li:post:"
	GeoUrnPrefix                 string = "urn:li:geo:"
	ConversionUrnPrefix          string = "urn:lla:llaPartnerConversion:"
	countDefault                 uint   = 10
	maxUrnsPerCall               uint   = 50
)

type Service struct {
	clientId      string
	oAuth2Service *oauth2.Service
	errorResponse *ErrorResponse
}

type ServiceConfig struct {
	ClientId     string
	ClientSecret string
	TokenSource  tokensource.TokenSource
	RedirectUrl  *string
}

// NewService return new instance of LinkedIn struct
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

	return &Service{
		clientId:      serviceConfig.ClientId,
		oAuth2Service: oAuth2Service,
	}, nil
}

func (service *Service) versionedHttpRequest(requestConfig *go_http.RequestConfig, linkedInVersion *string) (*http.Request, *http.Response, *errortools.Error) {
	headers := requestConfig.NonDefaultHeaders
	if headers == nil {
		headers = &http.Header{}
	}
	version := defaultLinkedInVersion
	if linkedInVersion != nil {
		version = *linkedInVersion
	}
	(*headers).Set(linkedInVersionHeader, version)

	requestConfig.NonDefaultHeaders = headers

	// add error model
	service.errorResponse = &ErrorResponse{}
	requestConfig.ErrorModel = service.errorResponse

	request, response, e := service.oAuth2Service.HttpRequest(requestConfig)
	if e != nil {
		if service.errorResponse.Message != "" {
			e.SetMessage(service.errorResponse.Message)
		}
	}

	return request, response, e
}

func (service *Service) urlV2(path string) string {
	return fmt.Sprintf("%s/v2/%s", apiUrl, path)
}

func (service *Service) urlRest(path string) string {
	return fmt.Sprintf("%s/%s", apiUrlRest, path)
}

func (service *Service) FromUrn(prefix string, urn string) int64 {
	id, err := strconv.ParseInt(strings.TrimPrefix(urn, prefix), 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func (service *Service) AuthorizeUrl(scope string, accessType *string, prompt *string, state *string) string {
	return service.oAuth2Service.AuthorizeUrl(&scope, accessType, prompt, state)
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

type CreatedModified struct {
	Actor string `json:"actor"`
	Time  int64  `json:"time"`
}
