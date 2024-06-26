// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/NotificationSrv.proto

package NotificationSrv

import (
	fmt "fmt"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for NotificationSrv service

func NewNotificationSrvEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for NotificationSrv service

type NotificationSrvService interface {
	CreateRequesterNotification(ctx context.Context, in *CreateRequesterNotificationRequest, opts ...client.CallOption) (*CreateRequesterNotificationResponse, error)
	CreateValidatorNotification(ctx context.Context, in *CreateValidatorNotificationRequest, opts ...client.CallOption) (*CreateValidatorNotificationResponse, error)
	CreateBulkValidatorNotification(ctx context.Context, in *CreateBulkValidatorNotificationRequest, opts ...client.CallOption) (*CreateBulkValidatorNotificationResponse, error)
	GetRequesterNotifications(ctx context.Context, in *GetRequesterNotificationsRequest, opts ...client.CallOption) (*GetRequesterNotificationsResponse, error)
	GetValidatorNotifications(ctx context.Context, in *GetValidatorNotificationsRequest, opts ...client.CallOption) (*GetValidatorNotificationsResponse, error)
	GetRequesterUnreadNotificationsCount(ctx context.Context, in *GetRequesterUnreadNotificationsCountRequest, opts ...client.CallOption) (*GetRequesterUnreadNotificationsCountResponse, error)
	GetValidatorUnreadNotificationsCount(ctx context.Context, in *GetValidatorUnreadNotificationsCountRequest, opts ...client.CallOption) (*GetValidatorUnreadNotificationsCountResponse, error)
	MarkNotificationStatusAsRead(ctx context.Context, in *MarkNotificationStatusAsReadRequest, opts ...client.CallOption) (*MarkNotificationStatusAsReadResponse, error)
}

type notificationSrvService struct {
	c    client.Client
	name string
}

func NewNotificationSrvService(name string, c client.Client) NotificationSrvService {
	return &notificationSrvService{
		c:    c,
		name: name,
	}
}

func (c *notificationSrvService) CreateRequesterNotification(ctx context.Context, in *CreateRequesterNotificationRequest, opts ...client.CallOption) (*CreateRequesterNotificationResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.CreateRequesterNotification", in)
	out := new(CreateRequesterNotificationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationSrvService) CreateValidatorNotification(ctx context.Context, in *CreateValidatorNotificationRequest, opts ...client.CallOption) (*CreateValidatorNotificationResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.CreateValidatorNotification", in)
	out := new(CreateValidatorNotificationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationSrvService) CreateBulkValidatorNotification(ctx context.Context, in *CreateBulkValidatorNotificationRequest, opts ...client.CallOption) (*CreateBulkValidatorNotificationResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.CreateBulkValidatorNotification", in)
	out := new(CreateBulkValidatorNotificationResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationSrvService) GetRequesterNotifications(ctx context.Context, in *GetRequesterNotificationsRequest, opts ...client.CallOption) (*GetRequesterNotificationsResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.GetRequesterNotifications", in)
	out := new(GetRequesterNotificationsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationSrvService) GetValidatorNotifications(ctx context.Context, in *GetValidatorNotificationsRequest, opts ...client.CallOption) (*GetValidatorNotificationsResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.GetValidatorNotifications", in)
	out := new(GetValidatorNotificationsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationSrvService) GetRequesterUnreadNotificationsCount(ctx context.Context, in *GetRequesterUnreadNotificationsCountRequest, opts ...client.CallOption) (*GetRequesterUnreadNotificationsCountResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.GetRequesterUnreadNotificationsCount", in)
	out := new(GetRequesterUnreadNotificationsCountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationSrvService) GetValidatorUnreadNotificationsCount(ctx context.Context, in *GetValidatorUnreadNotificationsCountRequest, opts ...client.CallOption) (*GetValidatorUnreadNotificationsCountResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.GetValidatorUnreadNotificationsCount", in)
	out := new(GetValidatorUnreadNotificationsCountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationSrvService) MarkNotificationStatusAsRead(ctx context.Context, in *MarkNotificationStatusAsReadRequest, opts ...client.CallOption) (*MarkNotificationStatusAsReadResponse, error) {
	req := c.c.NewRequest(c.name, "NotificationSrv.MarkNotificationStatusAsRead", in)
	out := new(MarkNotificationStatusAsReadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for NotificationSrv service

type NotificationSrvHandler interface {
	CreateRequesterNotification(context.Context, *CreateRequesterNotificationRequest, *CreateRequesterNotificationResponse) error
	CreateValidatorNotification(context.Context, *CreateValidatorNotificationRequest, *CreateValidatorNotificationResponse) error
	CreateBulkValidatorNotification(context.Context, *CreateBulkValidatorNotificationRequest, *CreateBulkValidatorNotificationResponse) error
	GetRequesterNotifications(context.Context, *GetRequesterNotificationsRequest, *GetRequesterNotificationsResponse) error
	GetValidatorNotifications(context.Context, *GetValidatorNotificationsRequest, *GetValidatorNotificationsResponse) error
	GetRequesterUnreadNotificationsCount(context.Context, *GetRequesterUnreadNotificationsCountRequest, *GetRequesterUnreadNotificationsCountResponse) error
	GetValidatorUnreadNotificationsCount(context.Context, *GetValidatorUnreadNotificationsCountRequest, *GetValidatorUnreadNotificationsCountResponse) error
	MarkNotificationStatusAsRead(context.Context, *MarkNotificationStatusAsReadRequest, *MarkNotificationStatusAsReadResponse) error
}

func RegisterNotificationSrvHandler(s server.Server, hdlr NotificationSrvHandler, opts ...server.HandlerOption) error {
	type notificationSrv interface {
		CreateRequesterNotification(ctx context.Context, in *CreateRequesterNotificationRequest, out *CreateRequesterNotificationResponse) error
		CreateValidatorNotification(ctx context.Context, in *CreateValidatorNotificationRequest, out *CreateValidatorNotificationResponse) error
		CreateBulkValidatorNotification(ctx context.Context, in *CreateBulkValidatorNotificationRequest, out *CreateBulkValidatorNotificationResponse) error
		GetRequesterNotifications(ctx context.Context, in *GetRequesterNotificationsRequest, out *GetRequesterNotificationsResponse) error
		GetValidatorNotifications(ctx context.Context, in *GetValidatorNotificationsRequest, out *GetValidatorNotificationsResponse) error
		GetRequesterUnreadNotificationsCount(ctx context.Context, in *GetRequesterUnreadNotificationsCountRequest, out *GetRequesterUnreadNotificationsCountResponse) error
		GetValidatorUnreadNotificationsCount(ctx context.Context, in *GetValidatorUnreadNotificationsCountRequest, out *GetValidatorUnreadNotificationsCountResponse) error
		MarkNotificationStatusAsRead(ctx context.Context, in *MarkNotificationStatusAsReadRequest, out *MarkNotificationStatusAsReadResponse) error
	}
	type NotificationSrv struct {
		notificationSrv
	}
	h := &notificationSrvHandler{hdlr}
	return s.Handle(s.NewHandler(&NotificationSrv{h}, opts...))
}

type notificationSrvHandler struct {
	NotificationSrvHandler
}

func (h *notificationSrvHandler) CreateRequesterNotification(ctx context.Context, in *CreateRequesterNotificationRequest, out *CreateRequesterNotificationResponse) error {
	return h.NotificationSrvHandler.CreateRequesterNotification(ctx, in, out)
}

func (h *notificationSrvHandler) CreateValidatorNotification(ctx context.Context, in *CreateValidatorNotificationRequest, out *CreateValidatorNotificationResponse) error {
	return h.NotificationSrvHandler.CreateValidatorNotification(ctx, in, out)
}

func (h *notificationSrvHandler) CreateBulkValidatorNotification(ctx context.Context, in *CreateBulkValidatorNotificationRequest, out *CreateBulkValidatorNotificationResponse) error {
	return h.NotificationSrvHandler.CreateBulkValidatorNotification(ctx, in, out)
}

func (h *notificationSrvHandler) GetRequesterNotifications(ctx context.Context, in *GetRequesterNotificationsRequest, out *GetRequesterNotificationsResponse) error {
	return h.NotificationSrvHandler.GetRequesterNotifications(ctx, in, out)
}

func (h *notificationSrvHandler) GetValidatorNotifications(ctx context.Context, in *GetValidatorNotificationsRequest, out *GetValidatorNotificationsResponse) error {
	return h.NotificationSrvHandler.GetValidatorNotifications(ctx, in, out)
}

func (h *notificationSrvHandler) GetRequesterUnreadNotificationsCount(ctx context.Context, in *GetRequesterUnreadNotificationsCountRequest, out *GetRequesterUnreadNotificationsCountResponse) error {
	return h.NotificationSrvHandler.GetRequesterUnreadNotificationsCount(ctx, in, out)
}

func (h *notificationSrvHandler) GetValidatorUnreadNotificationsCount(ctx context.Context, in *GetValidatorUnreadNotificationsCountRequest, out *GetValidatorUnreadNotificationsCountResponse) error {
	return h.NotificationSrvHandler.GetValidatorUnreadNotificationsCount(ctx, in, out)
}

func (h *notificationSrvHandler) MarkNotificationStatusAsRead(ctx context.Context, in *MarkNotificationStatusAsReadRequest, out *MarkNotificationStatusAsReadResponse) error {
	return h.NotificationSrvHandler.MarkNotificationStatusAsRead(ctx, in, out)
}

// Api Endpoints for QuestionSrv service

func NewQuestionSrvEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for QuestionSrv service

type QuestionSrvService interface {
	CreateBulkValidatorQuestions(ctx context.Context, in *CreateBulkValidatorQuestionsRequest, opts ...client.CallOption) (*CreateBulkValidatorQuestionsResponse, error)
	SetValidatorQuestionAnswer(ctx context.Context, in *SetValidatorQuestionAnswerRequest, opts ...client.CallOption) (*SetValidatorQuestionAnswerResponse, error)
	GetQuestionExamByID(ctx context.Context, in *GetQuestionExamByIDRequest, opts ...client.CallOption) (*GetQuestionExamByIDResponse, error)
}

type questionSrvService struct {
	c    client.Client
	name string
}

func NewQuestionSrvService(name string, c client.Client) QuestionSrvService {
	return &questionSrvService{
		c:    c,
		name: name,
	}
}

func (c *questionSrvService) CreateBulkValidatorQuestions(ctx context.Context, in *CreateBulkValidatorQuestionsRequest, opts ...client.CallOption) (*CreateBulkValidatorQuestionsResponse, error) {
	req := c.c.NewRequest(c.name, "QuestionSrv.CreateBulkValidatorQuestions", in)
	out := new(CreateBulkValidatorQuestionsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionSrvService) SetValidatorQuestionAnswer(ctx context.Context, in *SetValidatorQuestionAnswerRequest, opts ...client.CallOption) (*SetValidatorQuestionAnswerResponse, error) {
	req := c.c.NewRequest(c.name, "QuestionSrv.SetValidatorQuestionAnswer", in)
	out := new(SetValidatorQuestionAnswerResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionSrvService) GetQuestionExamByID(ctx context.Context, in *GetQuestionExamByIDRequest, opts ...client.CallOption) (*GetQuestionExamByIDResponse, error) {
	req := c.c.NewRequest(c.name, "QuestionSrv.GetQuestionExamByID", in)
	out := new(GetQuestionExamByIDResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for QuestionSrv service

type QuestionSrvHandler interface {
	CreateBulkValidatorQuestions(context.Context, *CreateBulkValidatorQuestionsRequest, *CreateBulkValidatorQuestionsResponse) error
	SetValidatorQuestionAnswer(context.Context, *SetValidatorQuestionAnswerRequest, *SetValidatorQuestionAnswerResponse) error
	GetQuestionExamByID(context.Context, *GetQuestionExamByIDRequest, *GetQuestionExamByIDResponse) error
}

func RegisterQuestionSrvHandler(s server.Server, hdlr QuestionSrvHandler, opts ...server.HandlerOption) error {
	type questionSrv interface {
		CreateBulkValidatorQuestions(ctx context.Context, in *CreateBulkValidatorQuestionsRequest, out *CreateBulkValidatorQuestionsResponse) error
		SetValidatorQuestionAnswer(ctx context.Context, in *SetValidatorQuestionAnswerRequest, out *SetValidatorQuestionAnswerResponse) error
		GetQuestionExamByID(ctx context.Context, in *GetQuestionExamByIDRequest, out *GetQuestionExamByIDResponse) error
	}
	type QuestionSrv struct {
		questionSrv
	}
	h := &questionSrvHandler{hdlr}
	return s.Handle(s.NewHandler(&QuestionSrv{h}, opts...))
}

type questionSrvHandler struct {
	QuestionSrvHandler
}

func (h *questionSrvHandler) CreateBulkValidatorQuestions(ctx context.Context, in *CreateBulkValidatorQuestionsRequest, out *CreateBulkValidatorQuestionsResponse) error {
	return h.QuestionSrvHandler.CreateBulkValidatorQuestions(ctx, in, out)
}

func (h *questionSrvHandler) SetValidatorQuestionAnswer(ctx context.Context, in *SetValidatorQuestionAnswerRequest, out *SetValidatorQuestionAnswerResponse) error {
	return h.QuestionSrvHandler.SetValidatorQuestionAnswer(ctx, in, out)
}

func (h *questionSrvHandler) GetQuestionExamByID(ctx context.Context, in *GetQuestionExamByIDRequest, out *GetQuestionExamByIDResponse) error {
	return h.QuestionSrvHandler.GetQuestionExamByID(ctx, in, out)
}
