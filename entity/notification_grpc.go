package entity

type CreateRequesterNotificationRequest struct {
	From int64
	To   int64

	DetailMessage DetailMessage
}

type CreateValidatorNotificationRequest struct {
	From int64
	To   int64

	DetailMessage DetailMessage
}

type CreateBulkValidatorNotificationRequest struct {
	Bulk []CreateValidatorNotificationRequest
}

type GetRequesterNotificationsRequest struct {
	To     int64
	Limit  int64
	Offset int64
}

type GetValidatorNotificationsRequest struct {
	To     int64
	Limit  int64
	Offset int64
}

type GetRequesterUnreadNotificationsCountRequest struct {
	To int64
}

type GetValidatorUnreadNotificationsCountRequest struct {
	To int64
}

type MarkNotificationStatusAsReadRequest struct {
	ID string
}
