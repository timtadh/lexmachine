package inst

import (
	"fmt"
	"strings"
)

const (
	CHAR  = iota // CHAR instruction op code: match a byte in the range [X, Y] (inclusive)
	SPLIT        // SPLIT instruction op code: split jump to both X and Y
	JMP          // JMP instruction op code: jmp to X
	MATCH        // MATCH instruction op code: match the string
)

// Inst represents an NFA byte code instruction
type Inst struct {
	Op uint8
	X  uint32
	Y  uint32
}

// InstSlice is a list of NFA instructions
type InstSlice []*Inst

// New creates a new instruction
func New(op uint8, x, y uint32) *Inst {
	return &Inst{
		Op: op,
		X:  x,
		Y:  y,
	}
}

// String humanizes the byte code
func (self Inst) String() (s string) {
	switch self.Op {
	case CHAR:
		if self.X == self.Y {
			s = fmt.Sprintf("CHAR   %d (%q)", self.X, string([]byte{byte(self.X)}))
		} else {
			s = fmt.Sprintf("CHAR   %d (%q), %d (%q)", self.X, string([]byte{byte(self.X)}), self.Y, string([]byte{byte(self.Y)}))
		}
	case SPLIT:
		s = fmt.Sprintf("SPLIT  %v, %v", self.X, self.Y)
	case JMP:
		s = fmt.Sprintf("JMP    %v", self.X)
	case MATCH:
		s = "MATCH"
	}
	return
}

// Serialize outputs machine readable assembly
func (self Inst) Serialize() (s string) {
	switch self.Op {
	case CHAR:
		s = fmt.Sprintf("CHAR   %d, %d", self.X, self.Y)
	case SPLIT:
		s = fmt.Sprintf("SPLIT  %v, %v", self.X, self.Y)
	case JMP:
		s = fmt.Sprintf("JMP    %v", self.X)
	case MATCH:
		s = "MATCH"
	}
	return
}

// String humanizes the byte code
func (self InstSlice) String() (s string) {
	s = "{\n"
	for i, inst := range self {
		if inst == nil {
			continue
		}
		if i < 10 {
			s += fmt.Sprintf("    0%v %v\n", i, inst)
		} else {
			s += fmt.Sprintf("    %v %v\n", i, inst)
		}
	}
	s += "}"
	return
}

// Serialize outputs machine readable assembly
func (self InstSlice) Serialize() (s string) {
	lines := make([]string, 0, len(self))
	for i, inst := range self {
		lines = append(lines, fmt.Sprintf("%3d %s", i, inst.Serialize()))
	}
	return strings.Join(lines, "\n")
}
