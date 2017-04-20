package machines

import (
	"bytes"
	"fmt"
)

import (
	. "github.com/timtadh/lexmachine/inst"
	"github.com/timtadh/lexmachine/queue"
)

type UnconsumedInput struct {
	StartTC     int
	FailTC      int
	StartLine   int
	StartColumn int
	FailLine    int
	FailColumn  int
	Text        []byte
}

func (u *UnconsumedInput) Error() string {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	stc := min(u.StartTC, len(u.Text)-1)
	etc := min(u.FailTC, len(u.Text))
	return fmt.Sprintf("Lexer error: could not match text starting at %v:%v failing at %v:%v.\n\tunmatched text: '%v'",
		u.StartLine, u.StartColumn,
		u.FailLine, u.FailColumn,
		string(u.Text[stc:etc]),
	)
}

type Match struct {
	PC          int
	TC          int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Bytes       []byte
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
	for i := prev_tc + 1; i <= tc && i < len(text); i++ {
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

type Scanner func(int) (int, *Match, error, Scanner)

func LexerEngine(program InstSlice, text []byte) Scanner {
	done := false
	match_pc := -1
	match_tc := -1

	prev_tc := 0
	line := 1
	col := 1

	var scan Scanner
	var cqueue, nqueue *queue.Queue = queue.New(len(program)), queue.New(len(program))
	scan = func(tc int) (int, *Match, error, Scanner) {
		if done && tc == len(text) {
			return tc, nil, nil, nil
		}
		start_tc := tc
		if tc < match_tc {
			// we back-tracked so reset the last match_tc
			match_tc = -1
		} else if tc == match_tc {
			// the caller did not reset the tc, we are where we left
		} else if match_tc != -1 && tc > match_tc {
			// we skipped text
			match_tc = tc
		}
		cqueue.Clear()
		nqueue.Clear()
		cqueue.Push(0)
		for ; tc <= len(text); tc++ {
			if cqueue.Empty() {
				break
			}
			for !cqueue.Empty() {
				pc := cqueue.Pop()
				inst := program[pc]
				switch inst.Op {
				case CHAR:
					x := byte(inst.X)
					y := byte(inst.Y)
					if tc < len(text) && x <= text[tc] && text[tc] <= y {
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
				default:
					panic(fmt.Errorf("unexpected instruction %v", inst))
				}
			}
			cqueue, nqueue = nqueue, cqueue
			if cqueue.Empty() && match_pc > -1 {
				line, col = compute_lc(text, prev_tc, start_tc, line, col)
				e_line, e_col := compute_lc(text, start_tc, match_tc-1, line, col)
				match := &Match{
					PC:          match_pc,
					TC:          start_tc,
					StartLine:   line,
					StartColumn: col,
					EndLine:     e_line,
					EndColumn:   e_col,
					Bytes:       text[start_tc:match_tc],
				}
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
			sline, scol := compute_lc(text, 0, start_tc, 1, 1)
			fline, fcol := compute_lc(text, 0, tc, 1, 1)
			err := &UnconsumedInput{
				StartTC:     start_tc,
				FailTC:      tc,
				StartLine:   sline,
				StartColumn: scol,
				FailLine:    fline,
				FailColumn:  fcol,
				Text:        text,
			}
			return tc, nil, err, scan
		} else {
			return tc, nil, nil, nil
		}
	}
	return scan
}
