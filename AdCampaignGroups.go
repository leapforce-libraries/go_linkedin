package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdCampaignGroupsResponse struct {
	Paging   Paging            `json:"paging"`
	Elements []AdCampaignGroup `json:"elements"`
}

type AdCampaignGroup struct {
	RunSchedule          AdRunSchedule       `json:"runSchedule"`
	ChangeAuditStamps    AdChangeAuditStamps `json:"changeAuditStamps"`
	Name                 string              `json:"name"`
	Test                 bool                `json:"test"`
	ServingStatuses      []string            `json:"servingStatuses"`
	AllowedCampaignTypes []AdCampaignType    `json:"allowedCampaignTypes"`
	Backfilled           bool                `json:"backfilled"`
	ID                   int64               `json:"id"`
	Account              string              `json:"account"`
	Status               string              `json:"status"`
}

type AdCampaignGroupStatus string

const (
	AdCampaignGroupStatusActive    AdCampaignGroupStatus = "ACTIVE"
	AdCampaignGroupStatusArchived  AdCampaignGroupStatus = "ARCHIVED"
	AdCampaignGroupStatusCanceled  AdCampaignGroupStatus = "CANCELED"
	AdCampaignGroupStatusDraft     AdCampaignGroupStatus = "DRAFT"
	AdCampaignGroupStatusCompleted AdCampaignGroupStatus = "COMPLETED"
)

type SearchAdCampaignGroupsConfig struct {
	Account *[]int64
	ID      *[]int64
	Status  *[]AdCampaignGroupStatus
	Name    *[]string
	Test    *bool
	Start   *uint
	Count   *uint
}

func (service *Service) SearchAdCampaignGroups(config *SearchAdCampaignGroupsConfig) (*[]AdCampaignGroup, *errortools.Error) {
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

	adCampaignGroups := []AdCampaignGroup{}

	for {
		if start > 0 {
			values.Set("start", fmt.Sprintf("%v", start))
		}
		if count != nil {
			values.Set("count", fmt.Sprintf("%v", *count))
		}

		adCampaignGroupsResponse := AdCampaignGroupsResponse{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("adCampaignGroupsV2?%s", values.Encode())),
			ResponseModel: &adCampaignGroupsResponse,
		}
		_, _, e := service.oAuth2Service.Get(&requestConfig)
		if e != nil {
			return nil, e
		}

		if len(adCampaignGroupsResponse.Elements) == 0 {
			break
		}

		adCampaignGroups = append(adCampaignGroups, adCampaignGroupsResponse.Elements...)

		if config != nil {
			if config.Start != nil {
				break
			}
		}

		if count == nil {
			_count := uint(adCampaignGroupsResponse.Paging.Count)
			count = &_count
		}

		start += *count

		if uint(adCampaignGroupsResponse.Paging.Total) <= start {
			break
		}
	}

	return &adCampaignGroups, nil
}

func (service *Service) GetAdCampaignGroups(accountID int64) (*[]AdCampaignGroup, *errortools.Error) {
	account := []int64{accountID}

	campaigns, e := service.SearchAdCampaignGroups(&SearchAdCampaignGroupsConfig{
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
