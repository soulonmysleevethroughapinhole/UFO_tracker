package tags

type TagInd struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	TargetPost int64  `json:"targetpost"`
	TagContent string `json:"tagcontent"`
}

type Tag struct {
	ID         int64  `json:"id"`
	PostID     int64  `json:"postid"`
	TagContent string `json:"tagcontent"`
	VoteAmt    int64  `json:"voteamt"`
	HasVoted   bool   `json:"hasvoted"`
}
