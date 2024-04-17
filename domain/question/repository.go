package question

import (
	"context"
)

type Repository interface {
	CreateBulkQuestions(ctx context.Context, questions []Question) error

	UpdateQuestion(ctx context.Context, questionID string, question Question) error

	GetQuestionerValidatorsIDsByRequestID(ctx context.Context, requesterID string) ([]int64, error)
	GetQuestionerValidatorCountByRequestID(ctx context.Context, requestID string) (int64, error)
	// questionID shoud be returned
	GetValidatorQuestionsByRequestID(ctx context.Context, validatorID int64, requestID string, limit, offset int64) ([]string, error)

	GetQuestionByID(ctx context.Context, id string) (*Question, error)
	// GetValidatorQuestion(ctx context.Context, requestId string, validatorId int64) (*Question, error)
}
