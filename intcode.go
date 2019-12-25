package intgode

import (
	"fmt"
)

var defaultInstructionSet = map[opcodeName]opcode{
	add:                addOpcode,
	multiply:           multiplyOpcode,
	input:              inputOpcode,
	output:             outputOpcode,
	jumpIfTrue:         jumpIfTrueOpcode,
	jumpIfFalse:        jumpIfFalseOpcode,
	lessThan:           lessThanOpcode,
	equals:             equalsOpcode,
	relativeBaseOffset: relativeBaseOffsetOpcode,
	halt:               haltOpcode,
}

// The IntcodeProgram interface represents an executable intcode program
type IntcodeProgram interface {
	// Starts the execution of the intcode program
	Exec()
	// Returns the memory state
	Memory() map[int]int
	// Returns whether the program has halted or not
	Halted() bool
	// Returns the relative base state
	RelativeBase() int
	// Returns the communication channel of the intcode program for sending input and receiving output (IO).
	IO() chan int
}

type intcodeProgram struct {
	io                 chan int
	halted             bool
	memory             map[int]int
	instructionPointer int
	instructionSet     map[opcodeName]opcode
	opcode             opcodeName
	parameterModes     [3]parameterMode
	relativeBase       int

	IntcodeProgram
}

func (ip *intcodeProgram) Exec() {
	io := ip.IO()

	defer func() {
		close(io)
	}()

	for {
		ip.parseOpcode()

		opcode, ok := ip.instructionSet[ip.opcode]
		if !ok {
			panic(fmt.Errorf("Illegal opcode %d", ip.opcode))
		}
		opcode(ip)

		if ip.halted {
			break
		}
	}
}

func (ip *intcodeProgram) Memory() map[int]int {
	currentMemory := make(map[int]int, len(ip.memory))
	for key, value := range ip.memory {
		currentMemory[key] = value
	}
	return currentMemory
}

func (ip *intcodeProgram) Halted() bool {
	return ip.halted
}

func (ip *intcodeProgram) IO() chan int {
	return ip.io
}

func (ip *intcodeProgram) RelativeBase() int {
	return ip.relativeBase
}

// NewIntcodeProgram returns an intcode program which uses an instruction set
// containing all of the opcodes from Advent of Code 2019 intcode puzzles.
func NewIntcodeProgram(intcode []int) IntcodeProgram {
	memory := make(map[int]int, len(intcode))
	for i, v := range intcode {
		memory[i] = v
	}

	return &intcodeProgram{
		io:                 make(chan int),
		halted:             false,
		memory:             memory,
		instructionPointer: 0,
		instructionSet:     defaultInstructionSet,
		parameterModes:     [3]parameterMode{positionMode, positionMode, positionMode},
	}
}

func (ip *intcodeProgram) parseOpcode() {
	n := ip.memory[ip.instructionPointer]
	ip.opcode = opcodeName(n % 100)
	ip.parameterModes[0] = parameterMode((n % 1000) / 100)
	ip.parameterModes[1] = parameterMode((n % 10000) / 1000)
	ip.parameterModes[2] = parameterMode((n % 100000) / 10000)
}

func (ip *intcodeProgram) movePointer(offset int) {
	ip.instructionPointer += offset
}

func (ip *intcodeProgram) readAt(offset int) int {
	switch ip.parameterModes[offset-1] {
	case positionMode:
		return ip.memory[ip.memory[ip.instructionPointer+offset]]
	case immediateMode:
		return ip.memory[ip.instructionPointer+offset]
	case relativeMode:
		return ip.memory[ip.relativeBase+ip.memory[ip.instructionPointer+offset]]
	default:
		panic("Invalid parameter mode")
	}
}

func (ip *intcodeProgram) writeAt(offset, value int) {
	switch ip.parameterModes[offset-1] {
	case positionMode:
		ip.memory[ip.memory[ip.instructionPointer+offset]] = value
	case immediateMode:
		panic(fmt.Errorf("Invalid parameter mode %d for writing to memory", immediateMode))
	case relativeMode:
		ip.memory[ip.relativeBase+ip.memory[ip.instructionPointer+offset]] = value
	default:
		panic("Invalid parameter mode")
	}
}
