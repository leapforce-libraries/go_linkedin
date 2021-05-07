package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdCreativesResponse struct {
	Paging   Paging       `json:"paging"`
	Elements []AdCreative `json:"elements"`
}

type AdCreative struct {
	Campaign          string              `json:"campaign"`
	ProcessingState   string              `json:"processingState"`
	ChangeAuditStamps AdChangeAuditStamps `json:"changeAuditStamps"`
	Test              bool                `json:"test"`
	ID                int64               `json:"id"`
	Reference         string              `json:"reference"`
	Review            struct {
		ReviewStatus string `json:"reviewStatus"`
	} `json:"review"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Variables struct {
		ClickURI string          `json:"clickUri"`
		Data     json.RawMessage `json:"data"`
	} `json:"variables"`
	Version AdVersion `json:"version"`
}

type AdCreativeStatus string

const (
	AdCreativeStatusActive   AdCreativeStatus = "ACTIVE"
	AdCreativeStatusPaused   AdCreativeStatus = "PAUSED"
	AdCreativeStatusDraft    AdCreativeStatus = "DRAFT"
	AdCreativeStatusArchived AdCreativeStatus = "ARCHIVED"
	AdCreativeStatusCanceled AdCreativeStatus = "CANCELED"
)

type AdCreativeType string

const (
	AdCreativeTypeTextAd           AdCreativeType = "TEXT_AD"
	AdCreativeTypeSponsoredUpdates AdCreativeType = "SPONSORED_UPDATES"
	AdCreativeTypeSponsoredInmails AdCreativeType = "SPONSORED_INMAILS"
	AdCreativeTypeDynamic          AdCreativeType = "DYNAMIC"
)

type SearchAdCreativesConfig struct {
	Campaign  *[]int64
	ID        *[]int64
	Reference *[]string
	Status    *[]AdCreativeStatus
	Test      *bool
	Start     *uint
	Count     *uint
}

func (service *Service) SearchAdCreatives(config *SearchAdCreativesConfig) (*[]AdCreative, *errortools.Error) {
	var values url.Values = url.Values{}
	var start uint = 0
	var count *uint = nil

	values.Set("q", "search")

	if config != nil {
		if config.Campaign != nil {
			for i, campaign := range *config.Campaign {
				values.Set(fmt.Sprintf("search.campaign.values[%v]", i), fmt.Sprintf("urn:li:sponsoredCampaign:%v", campaign))
			}
		}
		if config.ID != nil {
			for i, id := range *config.ID {
				values.Set(fmt.Sprintf("search.id.values[%v]", i), fmt.Sprintf("%v", id))
			}
		}
		if config.Reference != nil {
			for i, reference := range *config.Reference {
				values.Set(fmt.Sprintf("search.reference.values[%v]", i), reference)
			}
		}
		if config.Status != nil {
			for i, status := range *config.Status {
				values.Set(fmt.Sprintf("search.status.values[%v]", i), string(status))
			}
		}
		if config.Test != nil {
			values.Set("search.test", fmt.Sprintf("%v", *config.Test))
		}
		if config.Start != nil {
			start = *config.Start
		}
		if config.Count != nil {
			start = *config.Count
		}
	}

	adCreatives := []AdCreative{}

	for {
		if start > 0 {
			values.Set("start", fmt.Sprintf("%v", start))
		}
		if count != nil {
			values.Set("count", fmt.Sprintf("%v", *count))
		}

		adCreativesResponse := AdCreativesResponse{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("adCreativesV2?%s", values.Encode())),
			ResponseModel: &adCreativesResponse,
		}
		_, _, e := service.oAuth2Service.Get(&requestConfig)
		if e != nil {
			return nil, e
		}

		if len(adCreativesResponse.Elements) == 0 {
			break
		}

		adCreatives = append(adCreatives, adCreativesResponse.Elements...)

		if config != nil {
			if config.Start != nil {
				break
			}
		}

		if count == nil {
			_count := uint(adCreativesResponse.Paging.Count)
			count = &_count
		}

		start += *count

		if uint(adCreativesResponse.Paging.Total) <= start {
			break
		}
	}

	return &adCreatives, nil
}

func (service *Service) GetAdCreatives(campaignID int64) (*[]AdCreative, *errortools.Error) {
	campaign := []int64{campaignID}

	creatives, e := service.SearchAdCreatives(&SearchAdCreativesConfig{
		Campaign: &campaign,
	})
	if e != nil {
		return nil, e
	}

	if creatives == nil {
		return nil, nil
	}

	return creatives, nil
}
