package domain

type Chat int64
type User int

type MessageMembership struct {
	Conversation Chat
	Sender       User
}

type Message struct {
	Recipient Chat

	Text string
}
