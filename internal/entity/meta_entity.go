package entity

import "time"

const (
	QuestionEditSummaryKey = "question.edit.summary"
	QuestionCloseReasonKey = "question.close.reason"
	AnswerEditSummaryKey   = "answer.edit.summary"
	TagEditSummaryKey      = "tag.edit.summary"
)

// Meta meta
type Meta struct {
	ID        int       `xorm:"not null pk autoincr INT(10) id"`
	CreatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP created TIMESTAMP created_at"`
	UpdatedAt time.Time `xorm:"not null default CURRENT_TIMESTAMP updated TIMESTAMP updated_at"`
	ObjectID  string    `xorm:"not null default 0 INDEX BIGINT(20) object_id"`
	Key       string    `xorm:"not null VARCHAR(100) key"`
	Value     string    `xorm:"not null MEDIUMTEXT value"`
}

// TableName meta table name
func (Meta) TableName() string {
	return "meta"
}
