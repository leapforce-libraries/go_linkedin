package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	types "github.com/Leapforce-nl/go_types"
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
	Description        string               `json:"description"`
	Title              string               `json:"title"`
	ShareMediaCategory string               `json:"shareMediaCategory"`
}
type ShareContentEntity struct {
	Description    string      `json:"description"`
	EntityLocation string      `json:"entityLocation"`
	Thumbnails     []Thumbnail `json:"thumbnails"`
}

type Thumbnail struct {
	ImageSpecificContent json.RawMessage `json:"imageSpecificContent"`
	ResolvedUrl          string          `json:"resolvedUrl"`
}

func (li *LinkedIn) GetShares(organisationID int, startDateUnix int64, endDateUnix int64) (*[]Share, error) {
	if li == nil {
		return nil, &types.ErrorString{"Shares pointer is nil"}
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

		urlString := fmt.Sprintf("%s/shares?%s", li.BaseURL(), values.Encode())
		//fmt.Println(urlString)

		sharesResponse := SharesResponse{}

		_, err := li.OAuth2().Get(urlString, &sharesResponse)
		if err != nil {
			return nil, err
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
