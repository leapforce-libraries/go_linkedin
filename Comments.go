package linkedin

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type CommentsResponse struct {
	Paging   Paging    `json:"paging"`
	Elements []Comment `json:"elements"`
}

type Comment struct {
	Actor         *string              `json:"actor,omitempty"`
	Agent         *string              `json:"agent,omitempty"`
	Created       *CreatedModified     `json:"created,omitempty"`
	Id            *string              `json:"id,omitempty"`
	LastModified  *CreatedModified     `json:"lastModified,omitempty"`
	ParentComment *string              `json:"parentComment,omitempty"`
	Message       *CommentMessage      `json:"message,omitempty"`
	Urn           *string              `json:"$URN,omitempty"`
	LikesSummary  *CommentLikesSummary `json:"likesSummary,omitempty"`
	Content       *[]CommentContent    `json:"content,omitempty"`
	Object        *string              `json:"object,omitempty"`
}

type CommentLikesSummary struct {
	SelectedLikes        []string `json:"selectedLikes"`
	AggregatedTotalLikes int64    `json:"aggregatedTotalLikes"`
	LikedByCurrentUser   bool     `json:"likedByCurrentUser"`
	TotalLikes           int64    `json:"totalLikes"`
}

type CommentMessage struct {
	Attributes *[]CommentAttribute `json:"attributes,omitempty"`
	Text       string              `json:"text"`
}

type CommentAttribute struct {
	Length int64           `json:"length"`
	Start  int64           `json:"start"`
	Value  json.RawMessage `json:"value"`
}

type CommentContent struct {
	Type   string               `json:"type"`
	Entity CommentContentEntity `json:"entity"`
	Url    string               `json:"url"`
}

type CommentContentEntity struct {
	DigitalmediaAsset string `json:"digitalmediaAsset"`
}

func (service *Service) GetComments(urn string) (*[]Comment, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	comments := []Comment{}

	url := service.url(fmt.Sprintf("socialActions/%s/comments", urn))

	for {
		commentsResponse := CommentsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           url,
			ResponseModel: &commentsResponse,
		}
		_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, comment := range commentsResponse.Elements {
			comments = append(comments, comment)
		}

		url = ""
		for _, l := range commentsResponse.Paging.Links {
			if l.Rel == "next" {
				url = fmt.Sprintf("%s%s", apiUrlWithoutVersion, l.Href)
				break
			}
		}
		if url == "" {
			break
		}
	}

	return &comments, nil
}

func (service *Service) CreateComment(urn string, comment *Comment) (*Comment, *http.Response, *errortools.Error) {
	if service == nil {
		return nil, nil, errortools.ErrorMessage("Service pointer is nil")
	}
	if comment == nil {
		return nil, nil, errortools.ErrorMessage("Comment pointer is nil")
	}

	var newComment Comment

	url := service.urlRest(fmt.Sprintf("socialActions/%s/comments", urn))

	var header = http.Header{}
	header.Set(linkedInVersionHeader, defaultLinkedInVersion)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPost,
		Url:               url,
		BodyModel:         comment,
		ResponseModel:     &newComment,
		NonDefaultHeaders: &header,
	}
	_, resp, e := service.oAuth2Service.HttpRequest(&requestConfig)
	if e != nil {
		return nil, resp, e
	}

	return &newComment, resp, nil
}
