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
	Actor         string              `json:"actor"`
	Created       CreatedModified     `json:"created"`
	Id            string              `json:"id"`
	LastModified  CreatedModified     `json:"lastModified"`
	ParentComment string              `json:"parentComment"`
	Message       CommentMessage      `json:"message"`
	Urn           string              `json:"$URN"`
	LikesSummary  CommentLikesSummary `json:"likesSummary"`
	Content       []CommentContent    `json:"content"`
	Object        string              `json:"object"`
}

type CommentLikesSummary struct {
	SelectedLikes        []string `json:"selectedLikes"`
	AggregatedTotalLikes int64    `json:"aggregatedTotalLikes"`
	LikedByCurrentUser   bool     `json:"likedByCurrentUser"`
	TotalLikes           int64    `json:"totalLikes"`
}

type CommentMessage struct {
	Attributes []CommentAttribute `json:"attributes"`
	Text       string             `json:"text"`
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
