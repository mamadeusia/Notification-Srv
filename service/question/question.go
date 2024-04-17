package question

import (
	"context"
	"errors"

	"github.com/mamadeusia/NotificationSrv/domain/question"
	"github.com/mamadeusia/NotificationSrv/entity"
	"go-micro.dev/v4/events"
	"go-micro.dev/v4/logger"
)

// NotificationService - factory object
type QuestionService struct {
	stream events.Stream
	repo   question.Repository
}

func NewQuestionService(rm question.Repository, stream events.Stream) *QuestionService {
	return &QuestionService{
		repo:   rm,
		stream: stream,
	}
}

func (n *QuestionService) CreateBulkValidatorQuestions(ctx context.Context, req entity.CreateBulkValidatorQuestionsRequest) error {
	var questions []question.Question
	for _, q := range req.Bulk {
		questions = append(questions, question.Question{
			BaseQuestion: &entity.BaseQuestion{
				RequestID:    q.RequestID,
				ValidatorID:  q.ValidatorID,
				RequesterID:  q.RequesterID,
				QuestionType: q.DetailQuestion.GetQuestionType(),
			},
			QuestionDetails: q.DetailQuestion,
		})
	}
	if err := n.repo.CreateBulkQuestions(ctx, questions); err != nil {
		return err
	}
	return nil
}

func (n *QuestionService) SetValidatorQuestionAnswer(ctx context.Context, req entity.SetValidatorQuestionAnswerRequest) error {
	if err := n.repo.UpdateQuestion(ctx, req.QuestionID, question.Question{
		BaseQuestion: &entity.BaseQuestion{
			ID: req.QuestionID,
		},
		QuestionDetails: req.DetailAnswerOfQuestion,
	}); err != nil {
		logger.Info("QUESTIONSERVICE: SetValidatorQuestionAnswer has failed with error %v", err)
		return err
	}
	return nil
}

func (n *QuestionService) GetQuestionExamByID(ctx context.Context, req entity.GetQuestionExamByIDRequest) (*question.Question, error) {
	question, err := n.repo.GetQuestionByID(ctx, req.QuestionID)
	if err != nil {
		logger.Info("QUESTIONSERVICE: GetQuestionExamByID has failed with error %v", err)
		return nil, err
	}
	return question, nil
}

func (n *QuestionService) GetQuestionerValidatorForRequestID(ctx context.Context, requestID string, minRequiredValidator int) (validatorIDs []int64, err error) {

	count, err := n.repo.GetQuestionerValidatorCountByRequestID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	if count < int64(minRequiredValidator) {
		return nil, errors.New("not enough validator")
	}

	validatorIDs, err = n.repo.GetQuestionerValidatorsIDsByRequestID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	return validatorIDs, nil
}

func (n *QuestionService) GetQuestionsByValidatorIDs(ctx context.Context, requestID string, validatorIDs []int64) ([]string, error) {
	var questionIds []string
	for _, id := range validatorIDs {
		questionids, err := n.repo.GetValidatorQuestionsByRequestID(ctx, id, requestID, 10, 0)
		if err != nil {
			logger.Info("QUESTIONSERVICE: GetQuestionsByValidatorIDs has failed with error %v", err)
			return nil, err
		}
		questionIds = append(questionIds, questionids...)

	}

	return questionIds, nil
}
