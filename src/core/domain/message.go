package domain

//todo: create a struct that has Chat and User
type Chat int64
type User int

type Message struct {
	Recipient Chat

	Text string
}
