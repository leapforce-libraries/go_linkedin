package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdCampaignsResponse struct {
	Paging   Paging       `json:"paging"`
	Elements []AdCampaign `json:"elements"`
}

type AdCampaign struct {
	Account                  string              `json:"account"`
	AssociatedEntity         string              `json:"associatedEntity"`
	AudienceExpansionEnabled bool                `json:"audienceExpansionEnabled"`
	CampaignGroup            string              `json:"campaignGroup"`
	ChangeAuditStamps        AdChangeAuditStamps `json:"changeAuditStamps"`
	CostType                 string              `json:"costType"`
	CreativeSelection        string              `json:"creativeSelection"`
	DailyBudget              AdBudget            `json:"dailyBudget"`
	Format                   string              `json:"format"`
	ID                       int64               `json:"id"`
	Locale                   AdLocale            `json:"locale"`
	Name                     string              `json:"name"`
	ObjectiveType            string              `json:"objectiveType"`
	OffsiteDeliveryEnabled   bool                `json:"offsiteDeliveryEnabled"`
	OffsitePreferences       struct {
		IABCategories struct {
			Exclude []string `json:"exclude"`
			Include []string `json:"include"`
		} `json:"iabCategories"`
		PublisherRestrictionFiles struct {
			Exclude []string `json:"exclude"`
		} `json:"publisherRestrictionFiles"`
	} `json:"offsitePreferences"`
	OptimizationTargetType string        `json:"optimizationTargetType"`
	PacingStrategy         string        `json:"pacingStrategy"`
	RunSchedule            AdRunSchedule `json:"runSchedule"`
	ServingStatuses        []string      `json:"servingStatuses"`
	Status                 string        `json:"status"`
	Targeting              struct {
		IncludedTargetingFacets struct {
			Employers        []string   `json:"employers"`
			Locations        []string   `json:"locations"`
			InterfaceLocales []AdLocale `json:"interfaceLocales"`
		} `json:"includedTargetingFacets"`
	} `json:"targeting"`
	TargetingCriteria json.RawMessage `json:"targetingCriteria"`
	Test              bool            `json:"test"`
	TotalBudget       AdBudget        `json:"totalBudget"`
	Type              string          `json:"type"`
	UnitCost          AdBudget        `json:"unitCost"`
	Version           AdVersion       `json:"version"`
}

type AdCampaignStatus string

const (
	AdCampaignStatusActive    AdCampaignStatus = "ACTIVE"
	AdCampaignStatusPaused    AdCampaignStatus = "PAUSED"
	AdCampaignStatusArchived  AdCampaignStatus = "ARCHIVED"
	AdCampaignStatusCompleted AdCampaignStatus = "COMPLETED"
	AdCampaignStatusCanceled  AdCampaignStatus = "CANCELED"
	AdCampaignStatusDraft     AdCampaignStatus = "DRAFT"
)

type AdCampaignType string

const (
	AdCampaignTypeTextAd           AdCampaignType = "TEXT_AD"
	AdCampaignTypeSponsoredUpdates AdCampaignType = "SPONSORED_UPDATES"
	AdCampaignTypeSponsoredInmails AdCampaignType = "SPONSORED_INMAILS"
	AdCampaignTypeDynamic          AdCampaignType = "DYNAMIC"
)

type SearchAdCampaignsConfig struct {
	Account          *[]int64
	CampaignGroup    *[]int64
	AssociatedEntity *[]string
	ID               *[]int64
	Status           *[]AdCampaignStatus
	Type             *[]AdCampaignType
	Name             *[]string
	Test             *bool
	Start            *uint
	Count            *uint
}

func (service *Service) SearchAdCampaigns(config *SearchAdCampaignsConfig) (*[]AdCampaign, *errortools.Error) {
	var values url.Values = url.Values{}
	var start uint = 0
	var count *uint = nil

	values.Set("q", "search")

	if config != nil {
		if config.Account != nil {
			for i, account := range *config.Account {
				values.Set(fmt.Sprintf("search.account.values[%v]", i), fmt.Sprintf("urn:li:sponsoredAccount:%v", account))
			}
		}
		if config.CampaignGroup != nil {
			for i, campaignGroup := range *config.CampaignGroup {
				values.Set(fmt.Sprintf("search.campaignGroup.values[%v]", i), fmt.Sprintf("urn:li:sponsoredCampaignGroup:%v", campaignGroup))
			}
		}
		if config.AssociatedEntity != nil {
			for i, associatedEntity := range *config.AssociatedEntity {
				values.Set(fmt.Sprintf("search.associatedEntity.values[%v]", i), associatedEntity)
			}
		}
		if config.ID != nil {
			for i, id := range *config.ID {
				values.Set(fmt.Sprintf("search.id.values[%v]", i), fmt.Sprintf("%v", id))
			}
		}
		if config.Status != nil {
			for i, status := range *config.Status {
				values.Set(fmt.Sprintf("search.status.values[%v]", i), string(status))
			}
		}
		if config.Type != nil {
			for i, _type := range *config.Type {
				values.Set(fmt.Sprintf("search.type.values[%v]", i), string(_type))
			}
		}
		if config.Name != nil {
			for i, name := range *config.Name {
				values.Set(fmt.Sprintf("search.name.values[%v]", i), name)
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

	adCampaigns := []AdCampaign{}

	for {
		if start > 0 {
			values.Set("start", fmt.Sprintf("%v", start))
		}
		if count != nil {
			values.Set("count", fmt.Sprintf("%v", *count))
		}

		adCampaignsResponse := AdCampaignsResponse{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("adCampaignsV2?%s", values.Encode())),
			ResponseModel: &adCampaignsResponse,
		}
		_, _, e := service.oAuth2Service.Get(&requestConfig)
		if e != nil {
			return nil, e
		}

		if len(adCampaignsResponse.Elements) == 0 {
			break
		}

		adCampaigns = append(adCampaigns, adCampaignsResponse.Elements...)

		if config != nil {
			if config.Start != nil {
				break
			}
		}

		if count == nil {
			_count := uint(adCampaignsResponse.Paging.Count)
			count = &_count
		}

		start += *count

		if uint(adCampaignsResponse.Paging.Total) <= start {
			break
		}
	}

	return &adCampaigns, nil
}

func (service *Service) GetAdCampaigns(accountID int64) (*[]AdCampaign, *errortools.Error) {
	account := []int64{accountID}

	campaigns, e := service.SearchAdCampaigns(&SearchAdCampaignsConfig{
		Account: &account,
	})
	if e != nil {
		return nil, e
	}

	if campaigns == nil {
		return nil, nil
	}

	return campaigns, nil
}
