package domain

type InstructionType string

const (
	LimitInstruction InstructionType = "limit"
	StopInstruction  InstructionType = "stop"
)

type Instruction struct {
	Date        int64
	Recipient   Chat
	User        User
	Instruction InstructionType
	Parameter   int
}
