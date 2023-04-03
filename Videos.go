package linkedin

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"io/ioutil"
	"net/http"
)

type InitializeUploadVideoRequest struct {
	Owner           string `json:"owner"`
	FileSizeBytes   *int64 `json:"fileSizeBytes,omitempty"`
	UploadCaptions  *bool  `json:"uploadCaptions,omitempty"`
	UploadThumbnail *bool  `json:"uploadThumbnail,omitempty"`
}

type InitializeUploadVideoResponse struct {
	Value struct {
		UploadUrlsExpiresAt int64                              `json:"uploadUrlsExpireAt"`
		Video               string                             `json:"video"`
		UploadInstructions  []InitializeUploadVideoInstruction `json:"uploadInstructions"`
		UploadToken         string                             `json:"uploadToken"`
	} `json:"value"`
}

type InitializeUploadVideoInstruction struct {
	UploadUrl string `json:"uploadUrl"`
	FirstByte int64  `json:"firstByte"`
	LastByte  int64  `json:"lastByte"`
}

func (service *Service) InitializeUploadVideo(req *InitializeUploadVideoRequest) (*InitializeUploadVideoResponse, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}
	if req == nil {
		return nil, errortools.ErrorMessage("InitializeUploadVideoRequest pointer is nil")
	}

	var initializeUploadVideoRequest = struct {
		*InitializeUploadVideoRequest `json:"initializeUploadRequest"`
	}{req}

	var initializeUploadVideoResponse InitializeUploadVideoResponse

	var header = http.Header{}
	header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)
	header.Set(linkedInVersionHeader, defaultLinkedInVersion)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPost,
		Url:               service.urlRest("videos?action=initializeUpload"),
		BodyModel:         initializeUploadVideoRequest,
		ResponseModel:     &initializeUploadVideoResponse,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}
	return &initializeUploadVideoResponse, nil
}

func (service *Service) UploadVideo(uploadInstructions *[]InitializeUploadVideoInstruction, videoUrl string) (*[]string, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	resp, err := http.Get(videoUrl)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	var header = http.Header{}
	header.Set("Content-Type", "application/octet-stream")

	var etags []string

	for _, uploadInstruction := range *uploadInstructions {
		b := bytes[uploadInstruction.FirstByte:uploadInstruction.LastByte]

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodPut,
			Url:               uploadInstruction.UploadUrl,
			BodyRaw:           &b,
			NonDefaultHeaders: &header,
		}
		_, resp, e := service.oAuth2Service.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		etag := resp.Header.Get("etag")
		if etag == "" {
			return nil, errortools.ErrorMessage("UploadVideo did not return etag header")
		}

		etags = append(etags, etag)
	}

	return &etags, nil
}

type FinalizeUploadVideoRequest struct {
	Video           string   `json:"video"`
	UploadToken     string   `json:"uploadToken"`
	UploadedPartIds []string `json:"uploadedPartIds"`
}

func (service *Service) FinalizeUploadVideo(finalizeUploadVideoRequest *FinalizeUploadVideoRequest) *errortools.Error {
	if service == nil {
		return errortools.ErrorMessage("Service pointer is nil")
	}

	var finalizeUploadVideoRequest_ = struct {
		FinalizeUploadVideoRequest `json:"finalizeUploadRequest"`
	}{*finalizeUploadVideoRequest}

	var header = http.Header{}
	header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)
	header.Set(linkedInVersionHeader, defaultLinkedInVersion)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPost,
		Url:               service.urlRest("videos?action=finalizeUpload"),
		BodyModel:         finalizeUploadVideoRequest_,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)

	return e
}
