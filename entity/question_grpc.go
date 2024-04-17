package entity

type CreateValidatorQuestionRequest struct {
	RequestID   string
	ValidatorID int64
	RequesterID int64

	DetailQuestion DetailQuestion
}
type CreateBulkValidatorQuestionsRequest struct {
	Bulk []CreateValidatorQuestionRequest
}
type SetValidatorQuestionAnswerRequest struct {
	QuestionID string
	// we need to talk about it and checking of the validaty or it .
	DetailAnswerOfQuestion DetailQuestion
}
type GetQuestionExamByIDRequest struct {
	QuestionID string
}
