package entity

type QuestionConfirmedReq struct {
	RequestID            string
	RequesterID          int64
	MinRequiredQuestions int64
}
