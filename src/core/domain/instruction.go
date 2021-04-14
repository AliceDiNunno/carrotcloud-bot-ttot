package domain

type InstructionType string

const (
	LimitInstruction InstructionType = "limit"
	StopInstruction  InstructionType = "stop"
)

type Instruction struct {
	Metadata    MessageMetadata
	Instruction InstructionType
	Parameter   int
}
