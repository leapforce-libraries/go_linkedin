package linkedin

import (
	"fmt"
	"net/url"

	types "github.com/Leapforce-nl/go_types"
)

type UGCPostsResponse struct {
	Paging   Paging    `json:"paging"`
	Elements []UGCPost `json:"elements"`
}

type UGCPost struct {
	Author string `json:"author"`
}

func (li *LinkedIn) GetUGCPosts(organisationID int) (*[]UGCPost, error) {
	if li == nil {
		return nil, &types.ErrorString{"UGCPosts pointer is nil"}
	}

	values := url.Values{}
	values.Set("q", "authors")
	values.Set("authors", fmt.Sprintf("List({urn:li:organization:%v})", organisationID))

	urlString := fmt.Sprintf("%s/ugcPosts?%s", apiURL, values.Encode())
	urlString = "https://api.linkedin.com/v2/ugcPosts?q=authors&authors=LIST(urn%3Ali%3Aorganization%3A28586605)&sortBy=LAST_MODIFIED"
	fmt.Println(urlString)

	followerStatsResponse := UGCPostsResponse{}

	_, err := li.OAuth2().Get(urlString, &followerStatsResponse)
	if err != nil {
		return nil, err
	}

	return &followerStatsResponse.Elements, nil
}
