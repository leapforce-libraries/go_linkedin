package linkedin

type Paging struct {
	Count int    `json:"count"`
	Start int    `json:"start"`
	Links []Link `json:"links"`
	Total int    `json:"total"`
}

type Link struct {
	Type string `json:"type"`
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
