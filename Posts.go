package linkedin

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Post struct {
	Author                    string             `json:"author,omitempty"`
	Commentary                string             `json:"commentary,omitempty"`
	Visibility                string             `json:"visibility,omitempty"`
	Distribution              PostDistribution   `json:"distribution,omitempty"`
	Content                   *PostContent       `json:"content,omitempty"`
	LifecycleState            string             `json:"lifecycleState,omitempty"`
	LifecycleStateInfo        LifecycleStateInfo `json:"lifecycleStateInfo,omitempty"`
	IsRepostDisabledByAuthor  bool               `json:"isRepostDisabledByAuthor,omitempty"`
	IsReshareDisabledByAuthor bool               `json:"isReshareDisabledByAuthor,omitempty"`
	LastModifiedAt            int64              `json:"lastModifiedAt,omitempty"`
	PublishedAt               int64              `json:"publishedAt,omitempty"`
	CreatedAt                 int64              `json:"createdAt,omitempty"`
	Id                        string             `json:"id,omitempty"`
}

type PostDistribution struct {
	FeedDistribution               string   `json:"feedDistribution"`
	TargetEntities                 []string `json:"targetEntities,omitempty"`
	ThirdPartyDistributionChannels []string `json:"thirdPartyDistributionChannels,omitempty"`
}

type PostContent struct {
	Media PostContentMedia `json:"media"`
}

type PostContentMedia struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

type LifecycleStateInfo struct {
	IsEditedByAuthor bool `json:"isEditedByAuthor"`
}

func (service *Service) CreatePost(post *Post) *errortools.Error {
	if service == nil {
		return errortools.ErrorMessage("Service pointer is nil")
	}
	if post == nil {
		return errortools.ErrorMessage("Post pointer is nil")
	}

	var header = http.Header{}
	header.Set(linkedInVersionHeader, defaultLinkedInVersion)

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPost,
		Url:               service.urlRest("posts"),
		BodyModel:         post,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)

	return e
}

type PostsByOwnerConfig struct {
	OrganizationId           int64
	IsDirectSponsoredContent *bool
	Fields                   *string
	StartDateUnix            int64
	EndDateUnix              int64
}

type PostsByOwnerResponse struct {
	Paging   Paging `json:"paging"`
	Elements []Post `json:"elements"`
}

func (service *Service) PostsByOwner(cfg *PostsByOwnerConfig) (*[]Post, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}
	if cfg == nil {
		return nil, errortools.ErrorMessage("GetPostsByOwnerConfig pointer is nil")
	}

	start := 0
	count := 50

	var posts []Post

	for {
		values := url.Values{}
		values.Set("q", "author")
		values.Set("author", fmt.Sprintf("urn:li:organization:%v", cfg.OrganizationId))
		if cfg.IsDirectSponsoredContent != nil {
			values.Set("isDsc", fmt.Sprintf("%v", *cfg.IsDirectSponsoredContent))
		}
		if cfg.Fields != nil {
			values.Set("fields", *cfg.Fields)
		}
		values.Set("start", strconv.Itoa(start))
		values.Set("count", strconv.Itoa(count))

		postsResponse := PostsByOwnerResponse{}

		var header = http.Header{}
		header.Set(linkedInVersionHeader, defaultLinkedInVersion)

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("posts?%s", values.Encode())),
			ResponseModel:     &postsResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		fmt.Println(len(postsResponse.Elements))

		for _, post := range postsResponse.Elements {

			if post.CreatedAt > cfg.EndDateUnix {
				//continue
			}

			if post.CreatedAt < cfg.StartDateUnix {
				//continue
			}

			posts = append(posts, post)
		}

		if !postsResponse.Paging.HasLink("next") {
			break
		}

		start += count
	}

	return &posts, nil
}

type PostsResponse struct {
	Results map[string]Post `json:"results"`
}

func (service *Service) Posts(urns []string) (*[]Post, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	var posts []Post

	// deduplicate urns
	var _urnsMap = make(map[string]bool)
	var _urns []string
	for _, urn := range urns {
		_, ok := _urnsMap[urn]
		if ok {
			continue
		}
		_urnsMap[urn] = true
		_urns = append(_urns, urn)
	}

	for {
		var _urnsBatch []string

		if len(_urns) > int(maxUrnsPerCall) {
			_urnsBatch = _urns[:maxUrnsPerCall]
			_urns = _urns[maxUrnsPerCall:]
		} else {
			_urnsBatch = _urns
			_urns = []string{}
		}

		postsResponse := PostsResponse{}

		var header = http.Header{}
		header.Set(linkedInVersionHeader, defaultLinkedInVersion)

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.url(fmt.Sprintf("posts?ids=List(%s)", strings.Join(_urnsBatch, ","))),
			ResponseModel:     &postsResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, post := range postsResponse.Results {
			posts = append(posts, post)
		}

		if uint(len(_urns)) <= maxUrnsPerCall {
			break
		} else {
			_urns = _urns[maxUrnsPerCall:]
		}
	}

	return &posts, nil
}
