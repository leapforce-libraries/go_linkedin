package linkedin

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type ConversionsResponse struct {
	Paging   Paging       `json:"paging"`
	Elements []Conversion `json:"elements"`
}

type Conversion struct {
	PostClickAttributionWindowSize   *int64                          `json:"postClickAttributionWindowSize"`
	ViewThroughAttributionWindowSize *int64                          `json:"viewThroughAttributionWindowSize"`
	Created                          *int64                          `json:"created"`
	ImagePixelTag                    *string                         `json:"imagePixelTag"`
	Type                             *string                         `json:"type"`
	Enabled                          *bool                           `json:"enabled"`
	AssociatedCampaigns              *[]ConversionAssociatedCampaign `json:"associatedCampaigns"`
	Campaigns                        *[]string                       `json:"campaigns"`
	Name                             *string                         `json:"name"`
	LastModified                     *int64                          `json:"lastModified"`
	Id                               *int64                          `json:"id"`
	AttributionType                  *string                         `json:"attributionType"`
	UrlRules                         *[]UrlRule                      `json:"urlRules"`
	Value                            *ConversionValue                `json:"value"`
	Account                          *string                         `json:"account"`
}

type ConversionAssociatedCampaign struct {
	AssociatedAt int64  `json:"associatedAt"`
	Campaign     string `json:"campaign"`
	Conversion   string `json:"conversion"`
}

type ConversionValue struct {
	CurrencyCode string `json:"currencyCode"`
	Amount       string `json:"amount"`
}

type UrlRule struct {
	MatchValue string `json:"matchValue"`
	Type       string `json:"type"`
}

type GetConversionsConfig struct {
	AccountId int64
	Start     *uint
	Count     *uint
}

func (service *Service) GetConversionsForAccount(config *GetConversionsConfig) (*[]Conversion, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("config must not be nil")
	}

	var values = url.Values{}
	var start uint = 0
	var count uint = countDefault

	if config.Start != nil {
		start = *config.Start
	}
	if config.Count != nil {
		count = *config.Count
	}

	values.Set("q", "account")
	values.Set("account", fmt.Sprintf("%s%v", AccountUrnPrefix, config.AccountId))
	values.Set("count", fmt.Sprintf("%v", count))

	var conversions []Conversion

	for {
		if start > 0 {
			values.Set("start", fmt.Sprintf("%v", start))
		}

		var header = http.Header{}
		header.Set(restliProtocolVersionHeader, defaultRestliProtocolVersion)
		//header.Set("X-RestLi-Method", "FINDER")

		var conversionsResponse ConversionsResponse

		requestConfig := go_http.RequestConfig{
			Method:            http.MethodGet,
			Url:               service.urlRest(fmt.Sprintf("conversions?%s", values.Encode())),
			ResponseModel:     &conversionsResponse,
			NonDefaultHeaders: &header,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		if len(conversionsResponse.Elements) == 0 {
			break
		}

		conversions = append(conversions, conversionsResponse.Elements...)

		if config != nil {
			if config.Start != nil {
				break
			}
		}

		if len(conversionsResponse.Elements) < int(count) {
			break
		}

		start += count
	}

	return &conversions, nil
}
