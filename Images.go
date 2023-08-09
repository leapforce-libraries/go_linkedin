package linkedin

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"io"
	"net/http"
)

type InitializeUploadImageRequest struct {
	Owner string `json:"owner"`
}

type InitializeUploadImageResponse struct {
	Value struct {
		UploadUrlExpiresAt int64  `json:"uploadUrlExpiresAt"`
		UploadUrl          string `json:"uploadUrl"`
		Image              string `json:"image"`
	} `json:"value"`
}

func (service *Service) InitializeUploadImage(owner string) (*InitializeUploadImageResponse, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var initializeUploadRequest = struct {
		InitializeUploadImageRequest `json:"initializeUploadRequest"`
	}{InitializeUploadImageRequest{owner}}

	var initializeUploadResponse InitializeUploadImageResponse

	var header = http.Header{}
	header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)
	header.Set(linkedInVersionHeader, defaultLinkedInVersion)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPost,
		Url:               service.urlRest("images?action=initializeUpload"),
		BodyModel:         initializeUploadRequest,
		ResponseModel:     &initializeUploadResponse,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}
	return &initializeUploadResponse, nil
}

func (service *Service) UploadImage(putUrl string, imageUrl string) *errortools.Error {
	if service == nil {
		return errortools.ErrorMessage("Service pointer is nil")
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errortools.ErrorMessage(err)
	}

	var header = http.Header{}
	header.Set("Content-Type", http.DetectContentType(bytes))

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPut,
		Url:               putUrl,
		BodyRaw:           &bytes,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)

	return e
}

type Image struct {
	Owner  string `json:"owner"`
	Status string `json:"status"`
	Id     string `json:"id"`
}

func (service *Service) GetImage(imageUrn string, fields string) (*Image, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var header = http.Header{}
	header.Set(linkedInVersionHeader, defaultLinkedInVersion)

	var image Image

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodGet,
		Url:               service.urlRest(fmt.Sprintf("images/%s?fields=%s", imageUrn, fields)),
		ResponseModel:     &image,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &image, nil
}
