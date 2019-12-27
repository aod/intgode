package intgode

import "testing"

func TestOpcodeNames(t *testing.T) {
	testCases := []struct {
		name string
		want opcodeName
		got  opcodeName
	}{
		{"Add", 1, add},
		{"Multiply", 2, multiply},
		{"Input", 3, input},
		{"Output", 4, output},
		{"JumpIfTrue", 5, jumpIfTrue},
		{"JumpIfFalse", 6, jumpIfFalse},
		{"LessThan", 7, lessThan},
		{"Equals", 8, equals},
		{"RelativeBaseOffset", 9, relativeBaseOffset},
		{"Halt", 99, halt},
	}
	for _, tC := range testCases {
		if tC.want != tC.got {
			t.Errorf("Invalid opcode number for %s. Want %d got %d", tC.name, tC.want, tC.got)
		}
	}
}

func TestAddInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{1, 0, 0, 0, 99})
	go program.Exec()
	<-program.Output()

	if program.Memory()[0] != 2 {
		t.FailNow()
	}
}

func TestMultiplyInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{2, 0, 0, 0, 99})
	go program.Exec()
	<-program.Output()

	if program.Memory()[0] != 4 {
		t.FailNow()
	}
}

func TestInputInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{3, 0, 99})
	go program.Exec()

	<-program.Output()
	program.Input() <- 5

	<-program.Output()

	if program.Memory()[0] != 5 {
		t.FailNow()
	}
}

func TestOutputInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{4, 0, 99})
	go program.Exec()

	out := <-program.Output()

	if out[0] != 4 {
		t.FailNow()
	}
}

func TestJumpIfTrueInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{5, 0, 0, 4, 0, 99})
	go program.Exec()

	<-program.Output()

	if program.Memory()[0] != 5 {
		t.FailNow()
	}
}

func TestJumpIfFalseInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{6, 2, 0, 4, 0, 99, 1, 0, 0, 0, 99})
	go program.Exec()
	<-program.Output()

	if program.Memory()[0] != 12 {
		t.FailNow()
	}
}

func TestLessThanInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{7, 0, 4, 0, 99})
	go program.Exec()
	<-program.Output()

	if program.Memory()[0] != 1 {
		t.FailNow()
	}
}

func TestEqualsInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{8, 1, 1, 0, 99})
	go program.Exec()
	<-program.Output()

	if program.Memory()[0] != 1 {
		t.FailNow()
	}
}

func TestRelativeBaseOffsetInstruction(t *testing.T) {
	program := NewIntcodeProgram([]int{9, 2, 99})
	go program.Exec()
	<-program.Output()

	if program.RelativeBase() != 99 {
		t.FailNow()
	}
}

func TestHalt(t *testing.T) {
	program := NewIntcodeProgram([]int{99, 1, 0, 0, 0})
	go program.Exec()
	<-program.Output()

	if !program.Halted() || program.Memory()[0] != 99 {
		t.FailNow()
	}
}
