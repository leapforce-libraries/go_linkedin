package linkedin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	go_types "github.com/leapforce-libraries/go_types"
)

type AdAnalyticsResponse struct {
	// no pagination supported, see: https://docs.microsoft.com/en-us/linkedin/marketing/integrations/ads-reporting/ads-reporting
	//Paging   Paging        `json:"paging"`
	Elements []AdAnalytics `json:"elements"`
}

type AdAnalytics struct {
	ActionClicks                             int64                  `json:"actionClicks"`
	AdUnitClicks                             int64                  `json:"adUnitClicks"`
	Clicks                                   int64                  `json:"clicks"`
	Comments                                 int64                  `json:"comments"`
	CompanyPageClicks                        int64                  `json:"companyPageClicks"`
	ConversionValueInLocalCurrency           go_types.Float64String `json:"conversionValueInLocalCurrency"`
	CostInLocalCurrency                      go_types.Float64String `json:"costInLocalCurrency"`
	CostInUsd                                go_types.Float64String `json:"costInUsd"`
	DateRange                                AdDateRange            `json:"dateRange"`
	ExternalWebsiteConversions               int64                  `json:"externalWebsiteConversions"`
	ExternalWebsitePostClickConversions      int64                  `json:"externalWebsitePostClickConversions"`
	ExternalWebsitePostViewConversions       int64                  `json:"externalWebsitePostViewConversions"`
	Follows                                  int64                  `json:"follows"`
	FullScreenPlays                          int64                  `json:"fullScreenPlays"`
	Impressions                              int64                  `json:"impressions"`
	LandingPageClicks                        int64                  `json:"landingPageClicks"`
	LeadGenerationMailContactInfoShares      int64                  `json:"leadGenerationMailContactInfoShares"`
	LeadGenerationMailInterestedClicks       int64                  `json:"leadGenerationMailInterestedClicks"`
	Likes                                    int64                  `json:"likes"`
	OneClickLeadFormOpens                    int64                  `json:"oneClickLeadFormOpens"`
	OneClickLeads                            int64                  `json:"oneClickLeads"`
	Opens                                    int64                  `json:"opens"`
	OtherEngagements                         int64                  `json:"otherEngagements"`
	Pivot                                    string                 `json:"pivot"`
	PivotValue                               string                 `json:"pivotValue"`
	PivotValues                              []string               `json:"pivotValues"`
	Reactions                                int64                  `json:"reactions"`
	Sends                                    int64                  `json:"sends"`
	Shares                                   int64                  `json:"shares"`
	TextUrlClicks                            int64                  `json:"textUrlClicks"`
	TotalEngagements                         int64                  `json:"totalEngagements"`
	VideoCompletions                         int64                  `json:"videoCompletions"`
	VideoFirstQuartileCompletions            int64                  `json:"videoFirstQuartileCompletions"`
	VideoMidpointCompletions                 int64                  `json:"videoMidpointCompletions"`
	VideoStarts                              int64                  `json:"videoStarts"`
	VideoThirdQuartileCompletions            int64                  `json:"videoThirdQuartileCompletions"`
	VideoViews                               int64                  `json:"videoViews"`
	ViralCardClicks                          int64                  `json:"viralCardClicks"`
	ViralCardImpressions                     int64                  `json:"viralCardImpressions"`
	ViralClicks                              int64                  `json:"viralClicks"`
	ViralCommentLikes                        int64                  `json:"viralCommentLikes"`
	ViralComments                            int64                  `json:"viralComments"`
	ViralCompanyPageClicks                   int64                  `json:"viralCompanyPageClicks"`
	ViralExternalWebsiteConversions          int64                  `json:"viralExternalWebsiteConversions"`
	ViralExternalWebsitePostClickConversions int64                  `json:"viralExternalWebsitePostClickConversions"`
	ViralExternalWebsitePostViewConversions  int64                  `json:"viralExternalWebsitePostViewConversions"`
	ViralFollows                             int64                  `json:"viralFollows"`
	ViralFullScreenPlays                     int64                  `json:"viralFullScreenPlays"`
	ViralImpressions                         int64                  `json:"viralImpressions"`
	ViralLandingPageClicks                   int64                  `json:"viralLandingPageClicks"`
	ViralLikes                               int64                  `json:"viralLikes"`
	ViralOneClickLeadFormOpens               int64                  `json:"viralOneClickLeadFormOpens"`
	ViralOneClickLeads                       int64                  `json:"viralOneClickLeads"`
	ViralOtherEngagements                    int64                  `json:"viralOtherEngagements"`
	ViralReactions                           int64                  `json:"viralReactions"`
	ViralShares                              int64                  `json:"viralShares"`
	ViralTotalEngagements                    int64                  `json:"viralTotalEngagements"`
	ViralVideoCompletions                    int64                  `json:"viralVideoCompletions"`
	ViralVideoFirstQuartileCompletions       int64                  `json:"viralVideoFirstQuartileCompletions"`
	ViralVideoMidpointCompletions            int64                  `json:"viralVideoMidpointCompletions"`
	ViralVideoStarts                         int64                  `json:"viralVideoStarts"`
	ViralVideoThirdQuartileCompletions       int64                  `json:"viralVideoThirdQuartileCompletions"`
	ViralVideoViews                          int64                  `json:"viralVideoViews"`
}

type AdAnalyticsPivot string

const (
	AdAnalyticsPivotCompany                   AdAnalyticsPivot = "COMPANY"
	AdAnalyticsPivotAccount                   AdAnalyticsPivot = "ACCOUNT"
	AdAnalyticsPivotShare                     AdAnalyticsPivot = "SHARE"
	AdAnalyticsPivotCampaign                  AdAnalyticsPivot = "CAMPAIGN"
	AdAnalyticsPivotCreative                  AdAnalyticsPivot = "CREATIVE"
	AdAnalyticsPivotCampaignGroup             AdAnalyticsPivot = "CAMPAIGN_GROUP"
	AdAnalyticsPivotConversion                AdAnalyticsPivot = "CONVERSION"
	AdAnalyticsPivotConversionNode            AdAnalyticsPivot = "CONVERSATION_NODE"
	AdAnalyticsPivotConversionNodeOptionIndex AdAnalyticsPivot = "CONVERSATION_NODE_OPTION_INDEX"
	AdAnalyticsPivotServingLocation           AdAnalyticsPivot = "SERVING_LOCATION"
	AdAnalyticsPivotCardIndex                 AdAnalyticsPivot = "CARD_INDEX"
	AdAnalyticsPivotMemberCompanySize         AdAnalyticsPivot = "MEMBER_COMPANY_SIZE"
	AdAnalyticsPivotMemberIndustry            AdAnalyticsPivot = "MEMBER_INDUSTRY"
	AdAnalyticsPivotMemberSeniority           AdAnalyticsPivot = "MEMBER_SENIORITY"
	AdAnalyticsPivotMemberJobTitle            AdAnalyticsPivot = "MEMBER_JOB_TITLE"
	AdAnalyticsPivotMemberJobFunction         AdAnalyticsPivot = "MEMBER_JOB_FUNCTION"
	AdAnalyticsPivotMemberCountryV2           AdAnalyticsPivot = "MEMBER_COUNTRY_V2"
	AdAnalyticsPivotMemberRegionV2            AdAnalyticsPivot = "MEMBER_REGION_V2"
	AdAnalyticsPivotMemberCompany             AdAnalyticsPivot = "MEMBER_COMPANY"
)

type GetAdAnalyticsConfig struct {
	Pivot           AdAnalyticsPivot
	DateRange       AdDateRange
	TimeGranularity TimeGranularity
	CampaignType    *AdCampaignType
	Shares          *[]string
	Campaigns       *[]string
	Creatives       *[]string
	CampaignGroups  *[]string
	Accounts        *[]string
	Companies       *[]string
	Start           *uint
	Count           *uint
	Fields          *[]string
}

func (service *Service) GetAdAnalytics(config *GetAdAnalyticsConfig) (*[]AdAnalytics, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("GetAdAnalyticsConfig must not be nil")
	}

	var values = url.Values{}
	var itemType string
	var items []string
	var itemsPerBatch = 20

	values.Set("q", "analytics")
	values.Set("pivot", string(config.Pivot))
	if config.DateRange.Start != nil {
		values.Set("dateRange.start.day", fmt.Sprintf("%v", config.DateRange.Start.Day))
		values.Set("dateRange.start.month", fmt.Sprintf("%v", config.DateRange.Start.Month))
		values.Set("dateRange.start.year", fmt.Sprintf("%v", config.DateRange.Start.Year))
	}
	if config.DateRange.End != nil {
		values.Set("dateRange.end.day", fmt.Sprintf("%v", config.DateRange.End.Day))
		values.Set("dateRange.end.month", fmt.Sprintf("%v", config.DateRange.End.Month))
		values.Set("dateRange.end.year", fmt.Sprintf("%v", config.DateRange.End.Year))
	}
	values.Set("timeGranularity", string(config.TimeGranularity))
	if config.CampaignType != nil {
		values.Set("campaignType", string(*config.CampaignType))
	}
	if config.Shares != nil {
		itemType = "shares"
		for _, share := range *config.Shares {
			items = append(items, share)
			//itemValues.Set(fmt.Sprintf("shares", i), share)
		}
	} else if config.Campaigns != nil {
		itemType = "campaigns"
		for _, campaign := range *config.Campaigns {
			items = append(items, campaign)
			//values.Set(fmt.Sprintf("campaigns[%v]", i), campaign)
		}
	} else if config.Creatives != nil {
		itemType = "creatives"
		for _, creative := range *config.Creatives {
			items = append(items, creative)
			//.Set(fmt.Sprintf("creatives[%v]", i), creative)
		}
	} else if config.CampaignGroups != nil {
		itemType = "campaignGroups"
		for _, campaignGroup := range *config.CampaignGroups {
			items = append(items, campaignGroup)
			//values.Set(fmt.Sprintf("campaignGroups[%v]", i), campaignGroup)
		}
	} else if config.Accounts != nil {
		itemType = "accounts"
		for _, account := range *config.Accounts {
			items = append(items, account)
			//values.Set(fmt.Sprintf("accounts[%v]", i), account)
		}
	} else if config.Companies != nil {
		itemType = "companies"
		for _, company := range *config.Companies {
			items = append(items, company)
			//values.Set(fmt.Sprintf("companies[%v]", i), company)
		}
	}
	if config.Fields != nil {
		values.Set("fields", strings.Join(*config.Fields, ","))
	}

	var adAnalytics []AdAnalytics

	for {
		var values_ = values

		for i, item := range items {
			if i == itemsPerBatch {
				break
			}

			values_.Set(fmt.Sprintf("%s[%v]", itemType, i), item)
		}

		adAnalyticsResponse := AdAnalyticsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlRest(fmt.Sprintf("adAnalytics?%s", values_.Encode())),
			ResponseModel: &adAnalyticsResponse,
		}
		_, _, e := service.versionedHttpRequest(&requestConfig, nil)
		if e != nil {
			return nil, e
		}

		adAnalytics = append(adAnalytics, adAnalyticsResponse.Elements...)

		if len(items) <= itemsPerBatch {
			break
		}

		items = items[itemsPerBatch:]
	}

	return &adAnalytics, nil
}
