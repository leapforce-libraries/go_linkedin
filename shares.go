package linkedin

import (
	"encoding/json"
	"fmt"
	"net/url"

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

func (li *LinkedIn) GetShares(organisationID int) (*[]Share, error) {
	if li == nil {
		return nil, &types.ErrorString{"Shares pointer is nil"}
	}

	values := url.Values{}
	values.Set("q", "owners")
	values.Set("owners", fmt.Sprintf("urn:li:organization:%v", organisationID))

	urlString := fmt.Sprintf("%s/shares?%s", apiURL, values.Encode())
	//fmt.Println(urlString)

	sharesResponse := SharesResponse{}

	_, err := li.OAuth2().Get(urlString, &sharesResponse)
	if err != nil {
		return nil, err
	}

	return &sharesResponse.Elements, nil
}
