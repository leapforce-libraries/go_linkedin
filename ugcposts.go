package linkedin

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
)

type UGCPostsResponse struct {
	Paging   Paging    `json:"paging"`
	Elements []UGCPost `json:"elements"`
}

type UGCPost struct {
	Author string `json:"author"`
}

func (li *LinkedIn) GetUGCPosts(organisationID int) (*[]UGCPost, *errortools.Error) {
	if li == nil {
		return nil, errortools.ErrorMessage("UGCPosts pointer is nil")
	}

	values := url.Values{}
	values.Set("q", "authors")
	values.Set("authors", fmt.Sprintf("List({urn:li:organization:%v})", organisationID))

	urlString := fmt.Sprintf("%s/ugcPosts?%s", li.BaseURL(), values.Encode())
	urlString = "https://api.linkedin.com/v2/ugcPosts?q=authors&authors=LIST(urn%3Ali%3Aorganization%3A28586605)&sortBy=LAST_MODIFIED"
	//fmt.Println(urlString)

	followerStatsResponse := UGCPostsResponse{}

	_, _, e := li.OAuth2().Get(urlString, &followerStatsResponse, nil)
	if e != nil {
		return nil, e
	}

	return &followerStatsResponse.Elements, nil
}
