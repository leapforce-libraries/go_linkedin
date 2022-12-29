package linkedin

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"io/ioutil"
	"net/http"
)

type RegisterUploadAssetRecipe string

const (
	RegisterUploadAssetRecipeFeedshareImage            RegisterUploadAssetRecipe = "urn:li:digitalmediaRecipe:feedshare-image"
	RegisterUploadAssetRecipeCompanyUpdateArticleImage RegisterUploadAssetRecipe = "urn:li:digitalmediaRecipe:companyUpdate-article-image"
	RegisterUploadAssetRecipeSsuCarouselCardImage      RegisterUploadAssetRecipe = "urn:li:digitalmediaRecipe:ssu-carousel-card-image"
	RegisterUploadAssetRecipeRightRailLogoImage        RegisterUploadAssetRecipe = "urn:li:digitalmediaRecipe:rightRail-logo-image"
	RegisterUploadAssetRecipeAdsImage                  RegisterUploadAssetRecipe = "urn:li:digitalmediaRecipe:ads-image"
)

type RegisterUploadAssetRequest struct {
	Owner                    string                      `json:"owner"`
	Recipes                  []RegisterUploadAssetRecipe `json:"recipes"`
	ServiceRelationships     []ServiceRelationship       `json:"serviceRelationships"`
	SupportedUploadMechanism *[]string                   `json:"supportedUploadMechanism,omitempty"`
	FileSize                 *int64                      `json:"fileSize,omitempty"`
}

type ServiceRelationship struct {
	Identifier       string `json:"identifier"`
	RelationshipType string `json:"relationshipType"`
}

type RegisterUploadAssetResponse struct {
	Value struct {
		UploadMechanism struct {
			MediaUploadHttpRequest *struct {
				UploadUrl string `json:"uploadUrl"`
				Headers   struct {
					MediaTypeFamily string `json:"media-type-family"`
				} `json:"headers"`
			} `json:"com.linkedin.digitalmedia.uploading.MediaUploadHttpRequest"`
			MultipartUpload *struct {
				PartUploadRequests []struct {
					Headers struct {
						ContentType string `json:"Content-Type"`
					} `json:"headers"`
					ByteRange struct {
						LastByte  int `json:"lastByte"`
						FirstByte int `json:"firstByte"`
					} `json:"byteRange"`
					Url          string `json:"url"`
					UrlExpiresAt int64  `json:"urlExpiresAt"`
				} `json:"partUploadRequests"`
				Metadata string `json:"metadata"`
			} `json:"com.linkedin.digitalmedia.uploading.MultipartUpload"`
		} `json:"uploadMechanism"`
		Asset         string `json:"asset"`
		MediaArtifact string `json:"mediaArtifact"`
	} `json:"value"`
}

type InitializeUploadAssetInstruction struct {
	UploadUrl string `json:"uploadUrl"`
	FirstByte int64  `json:"firstByte"`
	LastByte  int64  `json:"lastByte"`
}

func (service *Service) RegisterUploadAsset(req *RegisterUploadAssetRequest) (*RegisterUploadAssetResponse, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}
	if req == nil {
		return nil, errortools.ErrorMessage("InitializeUploadAssetRequest pointer is nil")
	}

	var registerUploadAssetRequest = struct {
		*RegisterUploadAssetRequest `json:"registerUploadRequest"`
	}{req}

	var registerUploadAssetResponse RegisterUploadAssetResponse

	var header = http.Header{}
	header.Set("X-Restli-Protocol-Version", "2.0.0")
	header.Set("LinkedIn-Version", "202209")

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPost,
		Url:               service.urlRest("assets?action=registerUpload"),
		BodyModel:         registerUploadAssetRequest,
		ResponseModel:     &registerUploadAssetResponse,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}
	return &registerUploadAssetResponse, nil
}

func (service *Service) UploadAsset(putUrl string, url string) (string, *errortools.Error) {
	if service == nil {
		return "", errortools.ErrorMessage("Service pointer is nil")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", errortools.ErrorMessage(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errortools.ErrorMessage(err)
	}

	var header = http.Header{}
	header.Set("Content-Type", http.DetectContentType(bytes))

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPut,
		Url:               putUrl,
		BodyRaw:           &bytes,
		NonDefaultHeaders: &header,
	}
	_, resp, e := service.oAuth2Service.HttpRequest(&requestConfig)

	etag := resp.Header.Get("etag")

	return etag, e
}

type CompleteMultipartUploadAssetRequest struct {
	MediaArtifact       string               `json:"mediaArtifact"`
	Metadata            string               `json:"metadata"`
	PartUploadResponses []PartUploadResponse `json:"partUploadResponses"`
}
type PartUploadResponse struct {
	Headers struct {
		ETag string `json:"ETag"`
	} `json:"headers"`
	HttpStatusCode int `json:"httpStatusCode"`
}

func (service *Service) CompleteMultipartUploadAsset(completeMultipartUploadAssetRequest *CompleteMultipartUploadAssetRequest) *errortools.Error {
	if service == nil {
		return errortools.ErrorMessage("Service pointer is nil")
	}

	var completeMultipartUploadAssetRequest_ = struct {
		CompleteMultipartUploadAssetRequest `json:"completeMultipartUploadRequest"`
	}{*completeMultipartUploadAssetRequest}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlRest("assets?action=completeMultiPartUpload"),
		BodyModel: completeMultipartUploadAssetRequest_,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)

	return e
}
