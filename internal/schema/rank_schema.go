package schema

// GetRankPersonalWithPageReq get rank list page request
type GetRankPersonalWithPageReq struct {
	// page
	Page int `validate:"omitempty,min=1" form:"page"`
	// page size
	PageSize int `validate:"omitempty,min=1" form:"page_size"`
	// username
	Username string `validate:"omitempty,gt=0,lte=100" form:"username"`
	// user id
	UserID string `json:"-"`
}

// GetRankPersonalWithPageResp rank response
type GetRankPersonalWithPageResp struct {
	// create time
	CreatedAt int64 `json:"created_at"`
	// object id
	ObjectID string `json:"object_id"`
	// question id
	QuestionID string `json:"question_id"`
	// answer id
	AnswerID string `json:"answer_id"`
	// object type
	ObjectType string `json:"object_type" enums:"question,answer,tag,comment"`
	// title
	Title string `json:"title"`
	// content
	Content string `json:"content"`
	// reputation
	Reputation int `json:"reputation"`
	// rank type
	RankType string `json:"rank_type"`
}
