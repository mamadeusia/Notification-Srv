package entity

// BaseQuestion - is an entity demonstrating validators' questions
type BaseQuestion struct {
	ID          string
	RequestID   string `bson:"request_id"`
	ValidatorID int64  `bson:"validator_id"`
	RequesterID int64  `bson:"requester_id"`

	QuestionType QuestionType `bson:"question_type"`
}

type QuestionType string

const (
	DescriptiveQuestion    QuestionType = "DescriptiveQuestion"
	MultipleChoiceQuestion QuestionType = "MultipleChoiceQuestion"
)
