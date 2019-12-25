package intgode

type opcodeName int

const (
	add opcodeName = iota + 1
	multiply
	input
	output
	jumpIfTrue
	jumpIfFalse
	lessThan
	equals
	relativeBaseOffset

	halt = 99
)

type opcode func(*intcodeProgram)

func addOpcode(ip *intcodeProgram) {
	ip.writeAt(3, ip.readAt(1)+ip.readAt(2))
	ip.movePointer(4)
}

func multiplyOpcode(ip *intcodeProgram) {
	ip.writeAt(3, ip.readAt(1)*ip.readAt(2))
	ip.movePointer(4)
}

func inputOpcode(ip *intcodeProgram) {
	ip.writeAt(1, <-ip.io)
	ip.movePointer(2)
}

func outputOpcode(ip *intcodeProgram) {
	ip.io <- ip.readAt(1)
	ip.movePointer(2)
}

func jumpIfTrueOpcode(ip *intcodeProgram) {
	if ip.readAt(1) != 0 {
		ip.instructionPointer = ip.readAt(2)
	} else {
		ip.movePointer(3)
	}
}

func jumpIfFalseOpcode(ip *intcodeProgram) {
	if ip.readAt(1) == 0 {
		ip.instructionPointer = ip.readAt(2)
	} else {
		ip.movePointer(3)
	}
}

func lessThanOpcode(ip *intcodeProgram) {
	if ip.readAt(1) < ip.readAt(2) {
		ip.writeAt(3, 1)
	} else {
		ip.writeAt(3, 0)
	}

	ip.movePointer(4)
}

func equalsOpcode(ip *intcodeProgram) {
	if ip.readAt(1) == ip.readAt(2) {
		ip.writeAt(3, 1)
	} else {
		ip.writeAt(3, 0)
	}

	ip.movePointer(4)
}

func relativeBaseOffsetOpcode(ip *intcodeProgram) {
	ip.relativeBase += ip.readAt(1)
	ip.movePointer(2)
}

func haltOpcode(ip *intcodeProgram) {
	ip.halted = true
}
