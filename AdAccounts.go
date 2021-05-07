package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type AdAccountsResponse struct {
	Paging   Paging      `json:"paging"`
	Elements []AdAccount `json:"elements"`
}

type AdAccount struct {
	Test                           bool                `json:"test"`
	ChangeAuditStamps              AdChangeAuditStamps `json:"changeAuditStamps"`
	Currency                       string              `json:"currency"`
	ID                             int64               `json:"id"`
	Name                           string              `json:"name"`
	NotifiedOnCampaignOptimization bool                `json:"notifiedOnCampaignOptimization"`
	NotifiedOnCreativeApproval     bool                `json:"notifiedOnCreativeApproval"`
	NotifiedOnCreativeRejection    bool                `json:"notifiedOnCreativeRejection"`
	NotifiedOnEndOfCampaign        bool                `json:"notifiedOnEndOfCampaign"`
	Reference                      string              `json:"reference"`
	ServingStatuses                []string            `json:"servingStatuses"`
	Status                         string              `json:"status"`
	TotalBudget                    AdBudget            `json:"totalBudget"`
	TotalBudgetEndsAt              int64               `json:"totalBudgetEndsAt"`
	Type                           string              `json:"type"`
	Version                        AdVersion           `json:"version"`
}

type AdAccountStatus string

const (
	AdAccountStatusDraft    AdAccountStatus = "DRAFT"
	AdAccountStatusCanceled AdAccountStatus = "CANCELED"
	AdAccountStatusActive   AdAccountStatus = "ACTIVE"
)

type AdAccountType string

const (
	AdAccountTypeBusiness   AdAccountType = "BUSINESS"
	AdAccountTypeEnterprise AdAccountType = "ENTERPRISE"
)

type SearchAdAccountsConfig struct {
	Status    *[]AdAccountStatus
	Reference *[]string
	Name      *[]string
	ID        *[]int64
	Type      *[]AdAccountType
	Test      *bool
	Start     *uint
	Count     *uint
}

func (service *Service) SearchAdAccounts(config *SearchAdAccountsConfig) (*[]AdAccount, *errortools.Error) {
	var values url.Values = url.Values{}
	var start uint = 0
	var count *uint = nil

	values.Set("q", "search")

	if config != nil {
		if config.Status != nil {
			for i, status := range *config.Status {
				values.Set(fmt.Sprintf("search.status.values[%v]", i), string(status))
			}
		}
		if config.Reference != nil {
			for i, reference := range *config.Reference {
				values.Set(fmt.Sprintf("search.reference.values[%v]", i), reference)
			}
		}
		if config.Name != nil {
			for i, name := range *config.Name {
				values.Set(fmt.Sprintf("search.name.values[%v]", i), name)
			}
		}
		if config.ID != nil {
			for i, id := range *config.ID {
				values.Set(fmt.Sprintf("search.id.values[%v]", i), fmt.Sprintf("%v", id))
			}
		}
		if config.Type != nil {
			for i, _type := range *config.Type {
				values.Set(fmt.Sprintf("search.type.values[%v]", i), string(_type))
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

	adAccounts := []AdAccount{}

	for {
		if start > 0 {
			values.Set("start", fmt.Sprintf("%v", start))
		}
		if count != nil {
			values.Set("count", fmt.Sprintf("%v", *count))
		}

		adAccountsResponse := AdAccountsResponse{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("adAccountsV2?%s", values.Encode())),
			ResponseModel: &adAccountsResponse,
		}
		_, _, e := service.oAuth2Service.Get(&requestConfig)
		if e != nil {
			return nil, e
		}

		if len(adAccountsResponse.Elements) == 0 {
			break
		}

		adAccounts = append(adAccounts, adAccountsResponse.Elements...)

		if config != nil {
			if config.Start != nil {
				break
			}
		}

		if count == nil {
			_count := uint(adAccountsResponse.Paging.Count)
			count = &_count
		}

		start += *count

		if uint(adAccountsResponse.Paging.Total) <= start {
			break
		}
	}

	return &adAccounts, nil
}

func (service *Service) GetAdAccount(accountID int64) (*AdAccount, *errortools.Error) {
	id := []int64{accountID}

	accounts, e := service.SearchAdAccounts(&SearchAdAccountsConfig{
		ID: &id,
	})
	if e != nil {
		return nil, e
	}

	if accounts == nil {
		return nil, errortools.ErrorMessage("Account not found.")
	}

	if len(*accounts) == 0 {
		return nil, errortools.ErrorMessage("Account not found.")
	}

	if len(*accounts) > 1 {
		return nil, errortools.ErrorMessage("Multiple accounts found.")
	}

	return &(*accounts)[0], nil
}
