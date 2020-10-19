package linkedin

type TotalShareStatistics struct {
	UniqueImpressionsCount *int     `json:"uniqueImpressionsCount"`
	ClickCount             *int     `json:"clickCount"`
	Engagement             *float64 `json:"engagement"`
	LikeCount              *int     `json:"likeCount"`
	CommentCount           *int     `json:"commentCount"`
	ShareCount             *int     `json:"shareCount"`
	CommentMentionsCount   *int     `json:"commentMentionsCount"`
	ImpressionCount        *int     `json:"impressionCount"`
	ShareMentionsCount     *int     `json:"shareMentionsCount"`
}
