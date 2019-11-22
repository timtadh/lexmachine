package stream_machines

import (
	"github.com/timtadh/lexmachine/machines"
	"github.com/timtadh/lexmachine/stream"
)

type Scanner func() (*machines.Match, error, Scanner)

// DFALexerEngine does the actual tokenization of the byte slice text using the
// DFA state machine. If the lexing process fails the Scanner will return
// an UnconsumedInput error.
func DFALexerEngine(startState, errorState int, trans machines.DFATrans, accepting machines.DFAAccepting, text stream.Stream) Scanner {
	var scan Scanner
	scan = func() (*machines.Match, error, Scanner) {
		if text.EOS() {
			return nil, nil, nil
		}
		buf := make([]stream.Character, 0, 10)
		matchID := -1
		matchLH := -1
		state := startState
		if match, has := accepting[state]; has {
			matchID = match
			matchLH = -1
		}
		if !text.Started() {
			if !text.Advance(1) {
				return nil, nil, nil
			}
		}
		for lh := 0; state != errorState; lh++ {
			c, has := text.Peek(lh)
			if !has {
				break
			}
			buf = append(buf, c)
			state = trans[state][c.Byte]
			if match, has := accepting[state]; has {
				matchID = match
				matchLH = lh
			}
		}
		if match, has := accepting[state]; has {
			matchID = match
			matchLH = len(buf) - 1
		}
		if matchLH == -1 && matchID > -1 {
			err := &machines.EmptyMatchError{
				MatchID: matchID,
				TC:      buf[0].TC,
				Line:    buf[0].Line,
				Column:  buf[0].Column,
			}
			return nil, err, scan
		} else if matchID > -1 && matchLH >= 0 {
			lexeme := make([]byte, 0, matchLH+1)
			for _, c := range buf[:matchLH+1] {
				lexeme = append(lexeme, c.Byte)
			}
			match := &machines.Match{
				PC:          matchID,
				TC:          buf[0].TC,
				StartLine:   buf[0].Line,
				StartColumn: buf[0].Column,
				EndLine:     buf[matchLH].Line,
				EndColumn:   buf[matchLH].Column,
				Bytes:       lexeme,
			}
			text.Advance(matchLH + 1)
			return match, nil, scan
		} else {
			lexeme := make([]byte, 0, len(buf))
			for _, c := range buf {
				lexeme = append(lexeme, c.Byte)
			}
			err := &machines.UnconsumedInput{
				StartTC:     buf[0].TC,
				FailTC:      buf[len(buf)-1].TC,
				StartLine:   buf[0].Line,
				StartColumn: buf[0].Column,
				FailLine:    buf[len(buf)-1].Line,
				FailColumn:  buf[len(buf)-1].Column,
				Text:        lexeme,
			}
			return nil, err, scan
		}
	}
	return scan
}
