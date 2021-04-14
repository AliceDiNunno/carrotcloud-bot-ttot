package domain

type Chat int64
type User int
type ID int
type Date int64

type MessageMetadata struct {
	Id           ID
	Conversation Chat
	Sender       User
	SentDate     Date
}

type Message struct {
	Metadata MessageMetadata

	Text string
}
