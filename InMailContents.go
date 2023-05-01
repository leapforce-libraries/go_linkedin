package linkedin

import (
	"fmt"
	"net/http"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type InMailContent struct {
	Account        *string `json:"account,omitempty"`
	CreatedAt      *int64  `json:"createdAt,omitempty"`
	CreatedBy      *string `json:"createdBy,omitempty"`
	CustomFooter   *string `json:"customFooter,omitempty"`
	HtmlBody       *string `json:"htmlBody,omitempty"`
	Id             *string `json:"id,omitempty"`
	LastModifiedAt *int64  `json:"lastModifiedAt,omitempty"`
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`
	Name           *string `json:"name,omitempty"`
	Sender         *string `json:"sender,omitempty"`
	SubContent     *struct {
		FormSubContent          *InMailFormSubContent          `json:"form"`
		GuidedRepliesSubContent *InMailGuidedRepliesSubContent `json:"guidedReplies"`
		StandardSubContent      *InMailRegularSubContent       `json:"regular"`
	} `json:"subContent,omitempty"`
	Subject *string `json:"subject,omitempty"`
}

type InMailFormSubContent struct {
	CallToActionText   string `json:"callToActionText"`
	RightRailAdPicture string `json:"rightRailAdPicture"`
}

type InMailGuidedRepliesSubContent struct {
	SponsoredConversation string `json:"sponsoredConversation"`
	RightRailAdPicture    string `json:"rightRailAdPicture"`
}

type InMailRegularSubContent struct {
	CallToActionText           string `json:"callToActionText"`
	CallToActionLandingPageUrl string `json:"callToActionLandingPageUrl"`
	RightRailAdPicture         string `json:"rightRailAdPicture"`
}

func (service *Service) GetInMailContent(id string) (*InMailContent, *errortools.Error) {
	inMailContent := InMailContent{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlRest(fmt.Sprintf("inMailContents/%s", id)),
		ResponseModel: &inMailContent,
	}
	_, _, e := service.versionedHttpRequest(&requestConfig, nil)
	if e != nil {
		return nil, e
	}

	return &inMailContent, nil
}
