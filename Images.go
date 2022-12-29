package linkedin

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"io/ioutil"
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
	header.Set("X-Restli-Protocol-Version", "2.0.0")
	header.Set("LinkedIn-Version", "202209")

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

	bytes, err := ioutil.ReadAll(resp.Body)
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
