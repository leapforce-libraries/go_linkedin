package linkedin

import (
	"fmt"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdInMailContent struct {
	Account           string              `json:"account"`
	ChangeAuditStamps AdChangeAuditStamps `json:"changeAuditStamps"`
	Editable          bool                `json:"editable"`
	HTMLBody          string              `json:"htmlBody"`
	ID                int64               `json:"id"`
	LegalText         struct {
		RawText string `json:"rawText"`
	} `json:"legalText"`
	Name       string         `json:"name"`
	Sender     AdInMailSender `json:"sender"`
	SubContent struct {
		FormSubContent          *AdInMailFormSubContent          `json:"com.linkedin.ads.AdInMailFormSubContent"`
		GuidedRepliesSubContent *AdInMailGuidedRepliesSubContent `json:"com.linkedin.ads.AdInMailGuidedRepliesSubContent"`
		StandardSubContent      *AdInMailStandardSubContent      `json:"com.linkedin.ads.AdInMailStandardSubContent"`
	} `json:"subContent"`
	Subject string `json:"subject"`
}

type AdInMailSender struct {
	DisplayName      string `json:"displayName"`
	DisplayPictureV2 string `json:"displayPictureV2"`
	From             string `json:"from"`
	SenderType       string `json:"senderType"`
}

type AdInMailFormSubContent struct {
	Action     string `json:"action"`
	ActionText string `json:"actionText"`
	AdUnitV2   string `json:"adUnitV2"`
}

type AdInMailGuidedRepliesSubContent struct {
	SponsoredConversation string `json:"sponsoredConversation"`
	RightRailAdPicture    string `json:"rightRailAdPicture"`
}

type AdInMailStandardSubContent struct {
	Action     string `json:"action"`
	ActionText string `json:"actionText"`
	AdUnitV2   string `json:"adUnitV2"`
}

func (service *Service) GetAdInMailContent(adInMailContentID int64) (*AdInMailContent, *errortools.Error) {
	adInMailContent := AdInMailContent{}

	requestConfig := go_http.RequestConfig{
		URL:           service.url(fmt.Sprintf("adInMailContentsV2/%v", adInMailContentID)),
		ResponseModel: &adInMailContent,
	}
	_, _, e := service.oAuth2Service.Get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &adInMailContent, nil
}
