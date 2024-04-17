package entity

// BaseMessage - is responsible for all kind of message that we have in our system
type BaseMessage struct {
	ID   string
	From int64 `bson:"from"`
	To   int64 `bson:"to"`

	MessageType   MessageType   `bson:"message_type"` //search on enum in mongo and the performance of the query
	MessageStatus MessageStatus `bson:"message_status"`
}

type MessageType string

const (
	Unknown            MessageType = "unkown"             //in case we face an error we will have an unkown type
	AdminAskQuestion   MessageType = "adminAskQuestion"   // admin will ask requester at the begining
	RequestRejected    MessageType = "requestRejected"    // admin will reject the case after validation
	RequestApproved    MessageType = "requestApproved"    // admin will approve the case after validation
	ValidatorQuestions MessageType = "validatorQuestions" //admin will inform requester that validators' questions are
	// ready and will have a button on this type message
	NearRequestFound            MessageType = "nearRequestFound"            //validator will be informed about requester near them
	ElectedAsValidator          MessageType = "electedAsValidator"          // validator will be informed that s/he is chosen as a final valdator
	RequesterRespondToValidator MessageType = "requesterRespondToValidator" //requester respose to validator questions
)

type MessageStatus string

const (
	UnRead MessageStatus = "unread"
	Read   MessageStatus = "read"
)
