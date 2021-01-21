package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	errortools "github.com/leapforce-libraries/go_errortools"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
)

type SharesResponse struct {
	Paging   Paging  `json:"paging"`
	Elements []Share `json:"elements"`
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

func (service *Service) GetShares(organisationID int, startDateUnix int64, endDateUnix int64) (*[]Share, *errortools.Error) {
	if service == nil {
		return nil, errortools.ErrorMessage("Shares pointer is nil")
	}

	start := 0
	count := 50
	doNext := true

	shares := []Share{}

	for doNext {
		values := url.Values{}
		values.Set("q", "owners")
		values.Set("owners", fmt.Sprintf("urn:li:organization:%v", organisationID))
		values.Set("sortBy", "CREATED")
		values.Set("start", strconv.Itoa(start))
		values.Set("count", strconv.Itoa(count))

		sharesResponse := SharesResponse{}

		requestConfig := oauth2.RequestConfig{
			URL:           service.url(fmt.Sprintf("shares?%s", values.Encode())),
			ResponseModel: &sharesResponse,
		}
		_, _, e := service.oAuth2.Get(&requestConfig)
		if e != nil {
			return nil, e
		}

		if len(sharesResponse.Elements) == 0 {
			doNext = false
			continue
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

		//fmt.Println("shares", len(shares))

		start += count
	}

	return &shares, nil
}
