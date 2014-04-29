package machines

import (
	"fmt"
	"bytes"
)

import (
	"github.com/timtadh/lexmachine/queue"
	. "github.com/timtadh/lexmachine/inst"
)

type Match struct {
	PC    int
	Bytes []byte
}

func (self *Match) Equals(other *Match) bool {
	if self == nil && other == nil {
		return true
	} else if self == nil {
		return false
	} else if other == nil {
		return false
	}
	return self.PC == other.PC && bytes.Equal(self.Bytes, other.Bytes)
}

func (self Match) String() string {
	return fmt.Sprintf("<Match %v '%v'>", self.PC, string(self.Bytes))
}

type Scanner func()(*Match, error, Scanner)

func LexerEngine(program InstSlice, text []byte) Scanner {
	var cqueue, nqueue *queue.Queue = queue.New(), queue.New()
	match_pc := -1
	match_tc := -1
	start_tc := 0
	cqueue.Push(0)
	tc := 0
	done := false

	var scan Scanner
	scan = func() (*Match, error, Scanner) {
		if done {
			return nil, nil, nil
		}
		for ; tc <= len(text); tc++ {
			for !cqueue.Empty() {
				pc := cqueue.Pop()
				inst := program[pc]
				switch inst.Op {
				case CHAR:
					x := byte(inst.X)
					y := byte(inst.Y)
					if tc < len(text) && x <= text[tc] && text[tc] <= y  {
						nqueue.Push(pc + 1)
					}
				case MATCH:
					if match_tc < tc {
						match_pc = int(pc)
						match_tc = tc
					} else if match_pc > int(pc) {
						match_pc = int(pc)
						match_tc = tc
					}
				case JMP:
					cqueue.Push(inst.X)
				case SPLIT:
					cqueue.Push(inst.X)
					cqueue.Push(inst.Y)
				}
			}
			cqueue, nqueue = nqueue, cqueue
			if cqueue.Empty() && match_pc != -1 {
				match := &Match{
					match_pc,
					text[start_tc:match_tc],
				}
				cqueue.Push(0)
				start_tc = tc
				match_pc = -1
				return match, nil, scan
			}
		}
		if match_tc != len(text) {
			done = true
			if match_tc == -1 {
				match_tc = 0
			}
			return nil, fmt.Errorf("Unconsumed text, %d, '%s'", match_tc, text[match_tc:]), scan
		} else {
			return nil, nil, nil
		}
	}
	return scan
}

