package entity

type DetailQuestion interface {
	GetQuestionDetails() (map[string]interface{}, error)
	GetQuestionType() QuestionType
}

type DescriptiveQuestionDetails struct {
	Question string `bson:"question"`
	Answer   string `bson:"answer"`
}

func (e *DescriptiveQuestionDetails) GetQuestionDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"question": e.Question,
		"answer":   e.Answer,
	}, nil
}
func (e *DescriptiveQuestionDetails) GetQuestionType() QuestionType {
	return DescriptiveQuestion
}

type MultipleChoiceQuestionDetails struct {
	Question    string   `bson:"question"`
	Choises     []string `bson:"choises"`
	AnswerIndex int      `bson:"answer_index"`
}

func (e *MultipleChoiceQuestionDetails) GetQuestionDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"question":     e.Question,
		"choises":      e.Choises,
		"answer_index": e.AnswerIndex,
	}, nil
}

func (e *MultipleChoiceQuestionDetails) GetQuestionType() QuestionType {
	return MultipleChoiceQuestion
}
