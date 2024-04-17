package entity

import (
	"time"
)

type DetailMessage interface {
	GetMessageDetails() (map[string]interface{}, error)
	GetMessageType() MessageType
}

// validator notification details part

type NearRequestFoundDetails struct {
	RequestID string `bson:"request_id"`
	FullName  string `bson:"full_name"`
}

func (e *NearRequestFoundDetails) GetMessageDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"request_id": e.RequestID,
		"full_name":  e.FullName,
	}, nil
}
func (e *NearRequestFoundDetails) GetMessageType() MessageType {
	return NearRequestFound
}

type ElectedAsValidatorDetails struct {
	RequesterFullName string `bson:"requester_full_name"`
}

func (e *ElectedAsValidatorDetails) GetMessageDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"requester_full_name": e.RequesterFullName,
	}, nil
}
func (e *ElectedAsValidatorDetails) GetMessageType() MessageType {
	return ElectedAsValidator
}

type RequesterRespondToValidatorDetails struct {
	RequestID      string          `bson:"request_id"`
	StoredMessages []StoredMessage `bson:"stored_messages"`
}

type StoredMessage struct {
	ChatID    int64 `bson:"chat_id"`
	MessageID int32 `bson:"message_id"`
}

func (e *RequesterRespondToValidatorDetails) GetMessageDetails() (map[string]interface{}, error) {
	// TODO :: implement this
	return map[string]interface{}{
		"request_id":      e.RequestID,
		"stored_messages": e.StoredMessages,
	}, nil
}

func (e *RequesterRespondToValidatorDetails) GetMessageType() MessageType {
	return RequesterRespondToValidator
}

// requester notification details part

type AdminAskQuestionDetails struct {
	Question string `bson:"question"`
}

// GetMessageDetails implements DetailMessage
func (a *AdminAskQuestionDetails) GetMessageDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"question": a.Question,
	}, nil
}

func (e *AdminAskQuestionDetails) GetMessageType() MessageType {
	return AdminAskQuestion
}

type RequestRejectedDetails struct {
	Reason string `bson:"reason"`
	Time   time.Time
}

// GetMessageDetails implements DetailMessage
func (r *RequestRejectedDetails) GetMessageDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"reason": r.Reason,
		"time":   r.Time,
	}, nil
}

func (e *RequestRejectedDetails) GetMessageType() MessageType {
	return RequestRejected
}

type RequestApprovedDetails struct {
	Time time.Time
}

func (r *RequestApprovedDetails) GetMessageDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"time": r.Time,
	}, nil
}

// GetMessageDetails implements DetailMessage
func (r *RequestApprovedDetails) GetMessageType() MessageType {
	return RequestApproved
}

type ValidatorQuestionsDetails struct {
	ValidatorQuestionIds []string
}

func (r *ValidatorQuestionsDetails) GetMessageDetails() (map[string]interface{}, error) {
	return map[string]interface{}{
		"validator_question_ids": r.ValidatorQuestionIds,
	}, nil
}

// GetMessageDetails implements DetailMessage
func (r *ValidatorQuestionsDetails) GetMessageType() MessageType {
	return ValidatorQuestions
}
