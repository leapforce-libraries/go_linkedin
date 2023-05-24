package linkedin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type CreativesResponse struct {
	Paging   Paging     `json:"paging"`
	Elements []Creative `json:"elements"`
}

type Creative struct {
	Account            *string          `json:"account,omitempty"`
	Campaign           *string          `json:"campaign,omitempty"`
	Content            *CreativeContent `json:"content,omitempty"`
	CreatedAt          *int64           `json:"createdAt,omitempty"`
	CreatedBy          *string          `json:"createdBy,omitempty"`
	Id                 *string          `json:"id,omitempty"`
	InlineContent      *InlineContent   `json:"inlineContent,omitempty"`
	IntendedStatus     *string          `json:"intendedStatus,omitempty"`
	IsServing          *bool            `json:"isServing,omitempty"`
	IsTest             *bool            `json:"isTest,omitempty"`
	LastModifiedAt     *int64           `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string          `json:"lastModifiedBy,omitempty"`
	Review             *CreativeReview  `json:"review,omitempty"`
	ServingHoldReasons *[]string        `json:"servingHoldReasons,omitempty"`
}

type CreativeContent struct {
	Reference string                    `json:"reference"`
	TextAd    *CreativeContentTextAd    `json:"textAd"`
	Jobs      *CreativeContentJobs      `json:"jobs"`
	Spotlight *CreativeContentSpotlight `json:"spotlight"`
	Follow    *CreativeContentFollow    `json:"follow"`
}

type CreativeContentTextAd struct {
	Image       string `json:"image"`
	Description string `json:"description"`
	Headline    string `json:"headline"`
	LandingPage string `json:"landingPage"`
}

type CreativeContentJobs struct {
	Logo                   string `json:"logo"`
	ShowMemberProfilePhoto bool   `json:"showMemberProfilePhoto"`
	OrganizationName       string `json:"organizationName"`
	Headline               struct {
		PreApproved string `json:"preApproved"`
	} `json:"headline"`
	ButtonLabel struct {
		PreApproved string `json:"preApproved"`
	} `json:"buttonLabel"`
}

type CreativeContentSpotlight struct {
	CallToAction           string `json:"callToAction"`
	Description            string `json:"description"`
	Headline               string `json:"headline"`
	LandingPage            string `json:"landingPage"`
	Logo                   string `json:"logo"`
	OrganizationName       string `json:"organizationName"`
	ShowMemberProfilePhoto bool   `json:"showMemberProfilePhoto"`
}

type CreativeContentFollow struct {
	OrganizationName string `json:"organizationName"`
	Logo             string `json:"logo"`
	Headline         struct {
		PreApproved string `json:"preApproved"`
	} `json:"headline"`
	Description struct {
		PreApproved string `json:"preApproved"`
	} `json:"description"`
	CallToAction           string `json:"callToAction"`
	ShowMemberProfilePhoto bool   `json:"showMemberProfilePhoto"`
}

type InlineContent struct {
	Post struct {
		AdContext struct {
			DscAdAccount string `json:"dscAdAccount"`
			DscStatus    string `json:"dscStatus"`
		} `json:"adContext"`
		Author                    string `json:"author"`
		Commentary                string `json:"commentary"`
		Visibility                string `json:"visibility"`
		LifecycleState            string `json:"lifecycleState"`
		IsReshareDisabledByAuthor bool   `json:"isReshareDisabledByAuthor"`
		Content                   struct {
			Media struct {
				Title string `json:"title"`
				Id    string `json:"id"`
			} `json:"media"`
		} `json:"content"`
	} `json:"post"`
}

type CreativeReview struct {
	Status           string   `json:"status"`
	RejectionReasons []string `json:"rejectionReasons"`
}

type SearchCreativesConfig struct {
	Accounts                                *[]string
	Campaigns                               *[]string
	ContentReferences                       *[]string
	Creatives                               *[]string
	IntendedStatuses                        *[]string
	IsTestAccount                           *bool
	IsTotalIncluded                         *bool
	LeadgenCreativeCallToActionDestinations *[]string
	SortOrder                               *string
	Start                                   *uint
	Count                                   *uint
}

func (service *Service) SearchCreatives(config *SearchCreativesConfig) (*[]Creative, *errortools.Error) {
	var params []string
	var start uint = 0
	var count uint = countDefault

	params = append(params, "q=criteria")

	if config != nil {
		if config.Accounts != nil {
			if len(*config.Accounts) > 0 {
				params = append(params, fmt.Sprintf("accounts=List(%s)", url.QueryEscape(strings.Join(*config.Accounts, ","))))
			}
		}
		if config.Campaigns != nil {
			if len(*config.Campaigns) > 0 {
				params = append(params, fmt.Sprintf("campaigns=List(%s)", url.QueryEscape(strings.Join(*config.Campaigns, ","))))
			}
		}
		if config.ContentReferences != nil {
			if len(*config.ContentReferences) > 0 {
				params = append(params, fmt.Sprintf("contentReferences=List(%s)", strings.Join(*config.ContentReferences, ",")))
			}
		}
		if config.Creatives != nil {
			if len(*config.Creatives) > 0 {
				params = append(params, fmt.Sprintf("creatives=List(%s)", strings.Join(*config.Creatives, ",")))
			}
		}
		if config.IntendedStatuses != nil {
			if len(*config.IntendedStatuses) > 0 {
				params = append(params, fmt.Sprintf("intendedStatuses=(value:List(%s))", strings.Join(*config.IntendedStatuses, ",")))
			}
		}
		if config.IsTestAccount != nil {
			params = append(params, fmt.Sprintf("isTestAccount=%v", *config.IsTestAccount))
		}
		if config.IsTotalIncluded != nil {
			params = append(params, fmt.Sprintf("isTotalIncluded=%v", *config.IsTotalIncluded))
		}
		if config.LeadgenCreativeCallToActionDestinations != nil {
			if len(*config.LeadgenCreativeCallToActionDestinations) > 0 {
				params = append(params, fmt.Sprintf("leadgenCreativeCallToActionDestinations=List(%s)", strings.Join(*config.LeadgenCreativeCallToActionDestinations, ",")))
			}
		}
		if config.SortOrder != nil {
			params = append(params, fmt.Sprintf("sortOrder=%s", *config.SortOrder))
		}
		if config.Start != nil {
			start = *config.Start
		}
		if config.Count != nil {
			count = *config.Count
		}
	}

	var creatives []Creative

	for {
		params_ := params
		if start > 0 {
			params_ = append(params_, fmt.Sprintf("start=%v", start))
		}
		params_ = append(params_, fmt.Sprintf("count=%v", count))

		creativesResponse := CreativesResponse{}

		var header = http.Header{}
		header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)
		header.Set("X-RestLi-Method", "FINDER")

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("creatives?%s", strings.Join(params_, "&"))),
			ResponseModel:     &creativesResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		if len(creativesResponse.Elements) == 0 {
			break
		}

		creatives = append(creatives, creativesResponse.Elements...)

		if config != nil {
			if config.Start != nil {
				break
			}
		}

		if len(creativesResponse.Elements) < int(count) {
			break
		}

		start += count
	}

	return &creatives, nil
}
