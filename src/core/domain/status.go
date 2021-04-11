package domain

type Status struct {
	Recipient          Chat
	Sender             User
	DidSucceed         bool
	AdditionnalDetails string
}
