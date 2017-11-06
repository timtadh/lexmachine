package machines

type DFATrans [][256]int
type DFAAccepting map[int]int

// DFALexerEngine does the actual tokenization of the byte slice text using the
// DFA state machine. If the lexing process fails the Scanner will return
// an UnconsumedInput error.
func DFALexerEngine(startState, errorState int, trans DFATrans, accepting DFAAccepting, text []byte) Scanner {
	done := false
	matchID := -1
	matchTC := -1

	prevTC := 0
	line := 1
	col := 1

	var scan Scanner
	scan = func(tc int) (int, *Match, error, Scanner) {
		if done && tc == len(text) {
			return tc, nil, nil, nil
		}
		startTC := tc
		if tc < matchTC {
			// we back-tracked so reset the last matchTC
			matchTC = -1
		} else if tc == matchTC {
			// the caller did not reset the tc, we are where we left
		} else if matchTC != -1 && tc > matchTC {
			// we skipped text
			matchTC = tc
		}
		state := startState
		for ; tc < len(text); tc++ {
			if match, has := accepting[state]; has {
				matchID = match
				matchTC = tc
			}
			state = trans[state][text[tc]]
			if state == errorState && matchID > -1 {
				line, col = computeLineCol(text, prevTC, startTC, line, col)
				eLine, eCol := computeLineCol(text, startTC, matchTC-1, line, col)
				match := &Match{
					PC:          matchID,
					TC:          startTC,
					StartLine:   line,
					StartColumn: col,
					EndLine:     eLine,
					EndColumn:   eCol,
					Bytes:       text[startTC:matchTC],
				}
				prevTC = startTC
				matchID = -1
				return tc, match, nil, scan
			}
		}
		if match, has := accepting[state]; has {
			matchID = match
			matchTC = tc
		}
		if matchTC != len(text) && startTC >= len(text) {
			// the user has moved us farther than the text. Assume that was
			// the intent and return EOF.
			return tc, nil, nil, nil
		} else if matchTC != len(text) {
			done = true
			if matchTC == -1 {
				matchTC = 0
			}
			sline, scol := computeLineCol(text, 0, startTC, 1, 1)
			fline, fcol := computeLineCol(text, 0, tc, 1, 1)
			err := &UnconsumedInput{
				StartTC:     startTC,
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
