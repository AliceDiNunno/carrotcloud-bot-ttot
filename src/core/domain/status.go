package domain

type Status struct {
	MetaData           MessageMetadata
	DidSucceed         bool
	AdditionnalDetails string
}
