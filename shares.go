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

type SharesByOwnerResponse struct {
	Paging   Paging  `json:"paging"`
	Elements []Share `json:"elements"`
}

type SharesResponse struct {
	Results map[string]Share `json:"results"`
}

type Share struct {
	Owner        string                     `json:"owner"`
	Activity     string                     `json:"activity"`
	Edited       bool                       `json:"edited"`
	Created      CreatedModified            `json:"created"`
	Text         ShareText                  `json:"text"`
	LastModified CreatedModified            `json:"lastModified"`
	ID           string                     `json:"id"`
	Distribution map[string]json.RawMessage `json:"distribution"`
	Content      *ShareContent              `json:"content"`
}

type CreatedModified struct {
	Actor string `json:"actor"`
	Time  int64  `json:"time"`
}

type ShareText struct {
	Annotations []Annotation `json:"annotations"`
	Text        string       `json:"text"`
}

type Annotation struct {
	Length int64  `json:"length"`
	Start  int64  `json:"start"`
	Entity string `json:"entity"`
}

type ShareContent struct {
	ContentEntities    []ShareContentEntity `json:"contentEntities"`
	Description        *string              `json:"description"`
	Title              *string              `json:"title"`
	LandingPageURL     *string              `json:"landingPageUrl"`
	ShareMediaCategory *string              `json:"shareMediaCategory"`
}

type ShareContentEntity struct {
	Description    *string     `json:"description"`
	EntityLocation *string     `json:"entityLocation"`
	Title          *string     `json:"title"`
	Thumbnails     []Thumbnail `json:"thumbnails"`
	Entity         *string     `json:"entity"`
}

type Thumbnail struct {
	ImageSpecificContent *ImageSpecificContent `json:"imageSpecificContent"`
	ResolvedUrl          string                `json:"resolvedUrl"`
}

type ImageSpecificContent struct {
	Size   *int `json:"size"`
	Width  *int `json:"width"`
	Height *int `json:"height"`
}

func (service *Service) GetSharesByOwner(organizationID int64, startDateUnix int64, endDateUnix int64) (*[]Share, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	start := 0
	count := 50
	doNext := true

	shares := []Share{}

	for doNext {
		values := url.Values{}
		values.Set("q", "owners")
		values.Set("owners", fmt.Sprintf("urn:li:organization:%v", organizationID))
		values.Set("sortBy", "CREATED")
		values.Set("start", strconv.Itoa(start))
		values.Set("count", strconv.Itoa(count))

		sharesResponse := SharesByOwnerResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("shares?%s", values.Encode())),
			ResponseModel: &sharesResponse,
		}
		_, _, e := service.oAuth2Service.HTTPRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, share := range sharesResponse.Elements {

			if share.Created.Time > endDateUnix {
				continue
			}

			if share.Created.Time < startDateUnix {
				doNext = false
				break
			}

			shares = append(shares, share)
		}

		if !sharesResponse.Paging.HasLink("next") {
			break
		}

		start += count
	}

	return &shares, nil
}

func (service *Service) GetShares(urns []string) (*[]Share, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Service pointer is nil")
	}

	shares := []Share{}

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
			if uint(i) == maxURNsPerCall {
				break
			}

			param := fmt.Sprintf("ids[%v]=%s", i, urn)

			if i > 0 {
				params = fmt.Sprintf("%s&%s", params, param)
			} else {
				params = param
			}

			i++
		}

		sharesResponse := SharesResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("shares?%s", params)),
			ResponseModel: &sharesResponse,
		}
		_, _, e := service.oAuth2Service.HTTPRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, share := range sharesResponse.Results {
			shares = append(shares, share)
		}

		if uint(len(_urns)) <= maxURNsPerCall {
			break
		} else {
			_urns = _urns[maxURNsPerCall:]
		}
	}

	return &shares, nil
}
