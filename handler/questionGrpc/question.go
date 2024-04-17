package questionGrpc

import (
	"context"

	"github.com/mamadeusia/NotificationSrv/entity"
	pb "github.com/mamadeusia/NotificationSrv/proto"
	"github.com/mamadeusia/NotificationSrv/service/question"
	"go-micro.dev/v4/logger"
)

type Handler struct {
	service *question.QuestionService
}

func New(srv *question.QuestionService) *Handler {
	return &Handler{
		service: srv,
	}
}

func (n *Handler) CreateBulkValidatorQuestions(ctx context.Context, req *pb.CreateBulkValidatorQuestionsRequest, res *pb.CreateBulkValidatorQuestionsResponse) error {
	//todo:: do we need any validation and checking here???
	var createBulkValidatorQuestionsRequest entity.CreateBulkValidatorQuestionsRequest
	for _, q := range req.Questions {
		createValidatorQuestionRequest := entity.CreateValidatorQuestionRequest{
			RequestID:   q.RequestID,
			ValidatorID: q.ValidatorID,
			RequesterID: q.ValidatorID,
		}
		if q.GetDescriptiveQuestionExamDetails() != nil {
			createValidatorQuestionRequest.DetailQuestion = &entity.DescriptiveQuestionDetails{
				Question: q.GetDescriptiveQuestionExamDetails().Question,
			}
		} else if q.GetMultipleChoiceQuestionExamDetails() != nil {
			createValidatorQuestionRequest.DetailQuestion = &entity.MultipleChoiceQuestionDetails{
				Question: q.GetMultipleChoiceQuestionExamDetails().Question,
				Choises:  q.GetMultipleChoiceQuestionExamDetails().Choices,
			}

		}

		createBulkValidatorQuestionsRequest.Bulk = append(createBulkValidatorQuestionsRequest.Bulk, createValidatorQuestionRequest)
	}
	if err := n.service.CreateBulkValidatorQuestions(ctx, createBulkValidatorQuestionsRequest); err != nil {
		logger.Info("HANDLER::CreateBulkValidatorQuestions,has failed with error %+v", err)
		return err
	}
	return nil
}
func (n *Handler) SetValidatorQuestionAnswer(ctx context.Context, req *pb.SetValidatorQuestionAnswerRequest, res *pb.SetValidatorQuestionAnswerResponse) error {

	var setValidatorQuestionAnswer entity.SetValidatorQuestionAnswerRequest

	if req.GetDescriptiveQuestionAnswerDetails() != nil {
		setValidatorQuestionAnswer.DetailAnswerOfQuestion = &entity.DescriptiveQuestionDetails{
			Question: req.QuestionID,
			Answer:   req.GetDescriptiveQuestionAnswerDetails().Answer,
		}

	} else if req.GetMultipleChoiceQuestionAnswerDetails() != nil {
		setValidatorQuestionAnswer.DetailAnswerOfQuestion = &entity.MultipleChoiceQuestionDetails{
			Question:    req.QuestionID,
			AnswerIndex: int(req.GetMultipleChoiceQuestionAnswerDetails().Index),
		}
	}
	if err := n.service.SetValidatorQuestionAnswer(ctx, setValidatorQuestionAnswer); err != nil {
		logger.Info("HANDLER::SetValidatorQuestionAnswer,has failed with error %+v", err)
		return err
	}

	return nil
}

func (n *Handler) GetQuestionExamByID(ctx context.Context, req *pb.GetQuestionExamByIDRequest, res *pb.GetQuestionExamByIDResponse) error {
	q, err := n.service.GetQuestionExamByID(ctx, entity.GetQuestionExamByIDRequest{
		QuestionID: req.Id,
	})
	if err != nil {
		logger.Info("HANDLER::GetQuestionExamByID,has failed with error %+v", err)
		return err
	}

	res.QuestionID = q.BaseQuestion.ID
	res.RequestID = q.BaseQuestion.RequestID
	res.RequesterID = q.BaseQuestion.RequesterID

	if q.BaseQuestion.QuestionType == entity.DescriptiveQuestion {
		qDestails, err := q.QuestionDetails.GetQuestionDetails()
		if err != nil {
			return err
		}

		question, err := GetValueMap[string]("question", qDestails)
		if err != nil {
			return err
		}
		res.MessageOneof = &pb.GetQuestionExamByIDResponse_DescriptiveQuestionExamDetails{
			DescriptiveQuestionExamDetails: &pb.DescriptiveQuestionExamDetails{
				Question: question,
			},
		}

	} else if q.BaseQuestion.QuestionType == entity.MultipleChoiceQuestion {
		qDestails, err := q.QuestionDetails.GetQuestionDetails()
		if err != nil {
			return err
		}

		question, err := GetValueMap[string]("question", qDestails)
		if err != nil {
			return err
		}
		choices, err := GetValueMap[[]string]("choices", qDestails)
		if err != nil {
			return err
		}
		res.MessageOneof = &pb.GetQuestionExamByIDResponse_MultipleChoiceQuestionExamDetails{
			MultipleChoiceQuestionExamDetails: &pb.MultipleChoiceQuestionExamDetails{
				Question: question,
				Choices:  choices,
			},
		}
	}

	return nil
}
