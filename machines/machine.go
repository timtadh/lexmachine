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
	TC    int
	Line  int
	Column int
	Bytes []byte
}

func compute_lc(text []byte, prev_tc, tc, line, col int) (int, int) {
	if tc < prev_tc {
		for i := prev_tc; i > tc; i-- {
			if text[i] == '\n' {
				line -= 1
			}
		}
		col = 0
		for i := tc; i >= 0; i-- {
			if text[i] == '\n' {
				break
			}
			col += 1
		}
		return line, col
	}
	for i := prev_tc+1; i <= tc; i++ {
		if text[i] == '\n' {
			col = 0
			line += 1
		} else {
			col += 1
		}
	}
	if prev_tc == tc && tc == 0 && tc < len(text) {
		if text[tc] == '\n' {
			line += 1
			col -= 1
		}
	}
	return line, col
}

func (self *Match) Equals(other *Match) bool {
	if self == nil && other == nil {
		return true
	} else if self == nil {
		return false
	} else if other == nil {
		return false
	}
	return self.PC == other.PC && 
			self.Line == other.Line &&
			self.Column == other.Column &&
			bytes.Equal(self.Bytes, other.Bytes)
}

func (self Match) String() string {
	return fmt.Sprintf("<Match %d %d (%d, %d) '%v'>", self.PC, self.TC, self.Line, self.Column, string(self.Bytes))
}

type Scanner func(int)(int, *Match, error, Scanner)

func LexerEngine(program InstSlice, text []byte) Scanner {
	var cqueue, nqueue *queue.Queue = queue.New(), queue.New()
	cqueue.Push(0)
	done := false
	match_pc := -1
	match_tc := -1

	prev_tc := 0
	line := 1
	col := 1

	var scan Scanner
	scan = func(tc int) (int, *Match, error, Scanner) {
		if done {
			return tc, nil, nil, nil
		}
		start_tc := tc
		if tc < match_tc {
			// we back-tracked so reset the last match_tc
			match_tc = -1
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
			if cqueue.Empty() && match_pc > -1 {
				line, col = compute_lc(text, prev_tc, start_tc, line, col)
				match := &Match{
					PC: match_pc,
					TC: start_tc,
					Line: line,
					Column: col,
					Bytes: text[start_tc:match_tc],
				}
				cqueue.Push(0)
				prev_tc = start_tc
				match_pc = -1
				return tc, match, nil, scan
			}
		}
		if match_tc != len(text) {
			done = true
			if match_tc == -1 {
				match_tc = 0
			}
			line, col = compute_lc(text, 0, match_tc, 1, 1)
			return tc, nil, fmt.Errorf("Unconsumed text, %d (%d, %d), '%s'", match_tc, line, col, text[match_tc:]), scan
		} else {
			return tc, nil, nil, nil
		}
	}
	return scan
}

