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
	StartLine  int
	StartColumn int
	EndLine  int
	EndColumn int
	Bytes []byte
}

func compute_lc(text []byte, prev_tc, tc, line, col int) (int, int) {
	if tc < 0 {
		return line, col
	}
	if tc < prev_tc {
		for i := prev_tc; i > tc && i > 0; i-- {
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
	for i := prev_tc+1; i <= tc && i < len(text); i++ {
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
			self.StartLine == other.StartLine &&
			self.StartColumn == other.StartColumn &&
			self.EndLine == other.EndLine &&
			self.EndColumn == other.EndColumn &&
			bytes.Equal(self.Bytes, other.Bytes)
}

func (self Match) String() string {
	return fmt.Sprintf("<Match %d %d (%d, %d)-(%d, %d) '%v'>", self.PC, self.TC, self.StartLine, self.StartColumn, self.EndLine, self.EndColumn, string(self.Bytes))
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
				case CHJMP:
					x := byte(inst.X)
					y := byte(inst.Y)
					if tc < len(text) && x <= text[tc] && text[tc] <= y  {
						nqueue.Push(pc + 1)
					} else {
						cqueue.Push(pc + 2)
					}
				}
			}
			cqueue, nqueue = nqueue, cqueue
			if cqueue.Empty() && match_pc > -1 {
				line, col = compute_lc(text, prev_tc, start_tc, line, col)
				e_line, e_col := compute_lc(text, start_tc, match_tc-1, line, col)
				match := &Match{
					PC: match_pc,
					TC: start_tc,
					StartLine: line,
					StartColumn: col,
					EndLine: e_line,
					EndColumn: e_col,
					Bytes: text[start_tc:match_tc],
				}
				cqueue.Push(0)
				prev_tc = start_tc
				match_pc = -1
				return tc, match, nil, scan
			}
		}
		if match_tc != len(text) && start_tc >= len(text) {
			// the user has moved us farther than the text. Assume that was
			// the intent and return EOF.
			return tc, nil, nil, nil
		} else if match_tc != len(text) {
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

func DFALexerEngine(program InstSlice, text []byte) Scanner {
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
		pc := uint32(0)
		forloop: for ; tc <= len(text) && int(pc) < len(program); {
			fmt.Println(pc, tc, len(text))
			inst := program[pc]
			fmt.Print(inst)

			if tc < len(text) {
				fmt.Println(" ", text[tc])
			} else {
				fmt.Println()
			}

			switch inst.Op {
			case CHAR:
				x := byte(inst.X)
				y := byte(inst.Y)
				if tc < len(text) && x <= text[tc] && text[tc] <= y  {
					pc += 1
					tc += 1
				} else {
					done = true
					// s := ""
					// if tc < len(text) {
						// s = string(text[tc])
					// } else {
						// s = "EOF"
					// }
					break forloop
					// return tc, nil, fmt.Errorf("(dfa) expected char %v, %d (%d, %d), '%s'", inst, match_tc, line, col,s), nil
				}
			case MATCH:
				if match_tc < tc {
					match_pc = int(pc)
					match_tc = tc
				} else if match_pc > int(pc) {
					match_pc = int(pc)
					match_tc = tc
				}
				fmt.Println("---------->", "match", inst, tc, pc, match_pc, match_tc)
				// pc += 1
				break forloop
			case JMP:
				pc = inst.X
			case SPLIT:
				panic(fmt.Errorf("You must supply a DFA you gave an NFA"))
			case CHJMP:
				x := byte(inst.X)
				y := byte(inst.Y)
				if tc < len(text) && x <= text[tc] && text[tc] <= y  {
					pc = pc + 1
					tc += 1
				} else {
					pc = pc + 2
				}
			}
		}
		if match_pc > -1 {
			line, col = compute_lc(text, prev_tc, start_tc, line, col)
			e_line, e_col := compute_lc(text, start_tc, match_tc-1, line, col)
			match := &Match{
				PC: match_pc,
				TC: start_tc,
				StartLine: line,
				StartColumn: col,
				EndLine: e_line,
				EndColumn: e_col,
				Bytes: text[start_tc:match_tc],
			}
			prev_tc = start_tc
			match_pc = -1
			return tc, match, nil, scan
		}
		if match_tc != len(text) && start_tc >= len(text) {
			// the user has moved us farther than the text. Assume that was
			// the intent and return EOF.
			return tc, nil, nil, nil
		} else if match_tc != len(text) {
			done = true
			if match_tc == -1 {
				match_tc = 0
			}
			line, col = compute_lc(text, 0, match_tc, 1, 1)
			return tc, nil, fmt.Errorf("(dfa) Unconsumed text, %d (%d, %d), '%s'", match_tc, line, col, text[match_tc:]), nil
		} else {
			return tc, nil, nil, nil
		}
	}
	return scan
}

