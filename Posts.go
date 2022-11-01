package linkedin

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type Post struct {
	Author                    string           `json:"author"`
	Commentary                string           `json:"commentary"`
	Visibility                string           `json:"visibility"`
	Distribution              PostDistribution `json:"distribution"`
	Content                   *PostContent     `json:"content,omitempty"`
	LifecycleState            string           `json:"lifecycleState"`
	IsReshareDisabledByAuthor bool             `json:"isReshareDisabledByAuthor"`
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

func (service *Service) CreatePost(post *Post) *errortools.Error {
	if service == nil {
		return errortools.ErrorMessage("Service pointer is nil")
	}
	if post == nil {
		return errortools.ErrorMessage("Post pointer is nil")
	}

	var header = http.Header{}
	header.Set("LinkedIn-Version", "202209")

	requestConfig := go_http.RequestConfig{
		Method:            http.MethodPost,
		Url:               service.urlRest("posts"),
		BodyModel:         post,
		NonDefaultHeaders: &header,
	}
	_, _, e := service.oAuth2Service.HttpRequest(&requestConfig)

	return e
}
