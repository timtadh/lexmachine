package inst

import (
	"fmt"
	"strings"
)

const (
	CHAR = iota
	SPLIT
	JMP
	MATCH
)

type Inst struct {
	Op uint8
	X  uint32
	Y  uint32
}

type InstSlice []*Inst

func New(op uint8, x, y uint32) *Inst {
	self := new(Inst)
	self.Op = op
	self.X = x
	self.Y = y
	return self
}

func (self Inst) String() (s string) {
	switch self.Op {
	case CHAR:
		if self.X == self.Y {
			s = fmt.Sprintf("CHAR   %d (%s)", self.X, string([]byte{byte(self.X)}))
		} else {
			s = fmt.Sprintf("CHAR   %d (%s), %d (%s)", self.X, string([]byte{byte(self.X)}), self.Y, string([]byte{byte(self.Y)}))
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

func (self InstSlice) Serialize() (s string) {
	lines := make([]string, 0, len(self))
	for i, inst := range self {
		lines = append(lines, fmt.Sprintf("%3d %s", i, inst.Serialize()))
	}
	return strings.Join(lines, "\n")
}
