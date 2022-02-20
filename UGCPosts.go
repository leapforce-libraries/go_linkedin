package linkedin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type UgcPostsByOwnerResponse struct {
	Paging   Paging    `json:"paging"`
	Elements []UgcPost `json:"elements"`
}

type UgcPostsResponse struct {
	Results map[string]UgcPost `json:"results"`
}

type UgcPost struct {
	Author                     string                     `json:"author"`
	Created                    CreatedModified            `json:"created"`
	FirstPublishedAt           int64                      `json:"firstPublishedAt"`
	ID                         string                     `json:"id"`
	LastModified               CreatedModified            `json:"lastModified"`
	LifecycleState             string                     `json:"lifecycleState"`
	SpecificContent            map[string]UgcShareContent `json:"specificContent"`
	VersionTag                 string                     `json:"versionTag"`
	Visibility                 json.RawMessage            `json:"visibility"`
	Distribution               json.RawMessage            `json:"distribution"`
	ContentCertificationRecord string                     `json:"contentCertificationRecord"`
}

type Distribution struct {
	ExternalDistributionChannels []string `json:"externalDistributionChannels"`
	DistributedViaFollowFeed     bool     `json:"distributedViaFollowFeed"`
	FeedDistribution             string   `json:"feedDistribution"`
}

type UgcShareContent struct {
	ShareCommentary    *ShareCommentary `json:"shareCommentary"`
	Media              []Media          `json:"media"`
	ShareFeatures      ShareFeatures    `json:"shareFeatures"`
	ShareMediaCategory string           `json:"shareMediaCategory"`
}

type ShareCommentary struct {
	InferredLocale string      `json:"inferredLocale"`
	Attributes     []Attribute `json:"attributes"`
	Text           string      `json:"text"`
}

type Attribute struct {
	Length int64          `json:"length"`
	Start  int64          `json:"start"`
	Value  AttributeValue `json:"value"`
}

type AttributeValue struct {
	CompanyAttributedEntity *CompanyAttributedEntity `json:"com.linkedin.common.CompanyAttributedEntity"`
	HashtagAttributedEntity *HashtagAttributedEntity `json:"com.linkedin.common.HashtagAttributedEntity"`
}

type CompanyAttributedEntity struct {
	Company string `json:"company"`
}

type HashtagAttributedEntity struct {
	Hashtag string `json:"hashtag"`
}

type Text struct {
	Text string `json:"text"`
}

type Media struct {
	Description Text           `json:"description"`
	OriginalUrl string         `json:"originalUrl"`
	Recipes     []string       `json:"recipes"`
	Media       string         `json:"media"`
	Title       Text           `json:"title"`
	Thumbnails  []UgcThumbnail `json:"thumbnails"`
	Status      string         `json:"status"`
}

type UgcThumbnail struct {
	Width  *int   `json:"width"`
	Url    string `json:"url"`
	Height *int   `json:"height"`
}

type ShareFeatures struct {
	Hashtags []string `json:"hashtags"`
}

func (service *Service) GetUgcPostsByOwner(organizationID int64, startDateUnix int64, endDateUnix int64) (*[]UgcPost, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	start := 0
	count := 50
	doNext := true

	ugcPosts := []UgcPost{}

	for doNext {
		values := url.Values{}
		values.Set("q", "authors")
		values.Set("sortBy", "CREATED")
		values.Set("start", strconv.Itoa(start))
		values.Set("count", strconv.Itoa(count))

		ugcPostsResponse := UgcPostsByOwnerResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("ugcPosts?%s&authors=List(%s)", values.Encode(), url.QueryEscape(fmt.Sprintf("urn:li:organization:%v", organizationID)))),
			ResponseModel: &ugcPostsResponse,
		}

		// add authentication header
		header := http.Header{}
		header.Set("X-Restli-Protocol-Version", "2.0.0")
		requestConfig.NonDefaultHeaders = &header

		_, _, e := service.oAuth2Service.HTTPRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, ugcPost := range ugcPostsResponse.Elements {

			if ugcPost.Created.Time > endDateUnix {
				continue
			}

			if ugcPost.Created.Time < startDateUnix {
				doNext = false
				break
			}

			ugcPosts = append(ugcPosts, ugcPost)
		}

		if !ugcPostsResponse.Paging.HasLink("next") {
			break
		}

		start += count
	}

	return &ugcPosts, nil
}

func (service *Service) GetUgcPosts(urns []string) (*[]UgcPost, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	ugcPosts := []UgcPost{}

	// deduplicate urns
	var _urnsMap map[string]bool = make(map[string]bool)
	_urns := []string{}
	for _, urn := range urns {
		_, ok := _urnsMap[urn]
		if ok {
			continue
		}
		_urnsMap[urn] = true
		_urns = append(_urns, urn)
	}

	for {
		params := ""

		for i, urn := range _urns {
			if uint(i) == maxUrnsPerCall {
				break
			}

			param := fmt.Sprintf("ids[%v]=%s", i, urn)

			if i > 0 {
				params = fmt.Sprintf("%s&%s", params, param)
			} else {
				params = param
			}
		}

		ugcPostsResponse := UgcPostsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("ugcPosts?%s", params)),
			ResponseModel: &ugcPostsResponse,
		}
		_, _, e := service.oAuth2Service.HTTPRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, ugcPost := range ugcPostsResponse.Results {
			ugcPosts = append(ugcPosts, ugcPost)
		}

		if uint(len(_urns)) <= maxUrnsPerCall {
			break
		} else {
			_urns = _urns[maxUrnsPerCall+1:]
		}
	}

	return &ugcPosts, nil
}

func (service *Service) GetUgcPost(urn string) (*UgcPost, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	ugcPost := UgcPost{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("ugcPosts/%s", url.QueryEscape(urn))),
		ResponseModel: &ugcPost,
	}

	// add authentication header
	header := http.Header{}
	header.Set("X-Restli-Protocol-Version", "2.0.0")
	requestConfig.NonDefaultHeaders = &header

	_, _, e := service.oAuth2Service.HTTPRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &ugcPost, nil
}
