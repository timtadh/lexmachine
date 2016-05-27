package machines

import (
	"fmt"
)

import ()

import (
	. "github.com/timtadh/lexmachine/inst"
)

func DFALexerEngine(program InstSlice, text []byte) Scanner {
	done := false
	line := 1
	col := 1
	prev_tc := 0
	match_tc := -1
	var scan Scanner
	scan = func(tc int) (int, *Match, error, Scanner) {
		if done || match_tc == len(text) {
			return tc, nil, nil, nil
		}
		start_tc := tc
		pc := 0
		loop: for ; tc <= len(text) && int(pc) < len(program); {
			inst := program[pc]
			fmt.Println(tc, len(text), inst)
			switch inst.Op {
			case CHAR:
				x := byte(inst.X)
				y := byte(inst.Y)
				if tc < len(text) && x <= text[tc] && text[tc] <= y  {
					pc += 1
					tc += 1
				} else {
					break loop
				}
			case MATCH:
				line, col = compute_lc(text, prev_tc, start_tc, line, col)
				e_line, e_col := compute_lc(text, start_tc, tc-1, line, col)
				match := &Match{
					PC: pc,
					TC: start_tc,
					StartLine: line,
					StartColumn: col,
					EndLine: e_line,
					EndColumn: e_col,
					Bytes: text[start_tc:tc],
				}
				match_tc = tc
				return tc, match, nil, scan
			case JMP:
				pc = int(inst.X)
			case SPLIT:
				panic(fmt.Errorf("You must supply a DFA you gave an NFA"))
			case CHJMP:
				x := byte(inst.X)
				y := byte(inst.Y)
				if tc < len(text) && x <= text[tc] && text[tc] <= y  {
					pc += 1
					tc += 1
				} else {
					pc += 2
				}
			}
		}
		done = true
		return tc, nil, fmt.Errorf("Unconsumed text, %d (%d, %d), '%s'", tc, line, col, text[tc:]), scan
	}
	return scan
}

