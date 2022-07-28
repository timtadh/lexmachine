package lexmachine

import (
	"fmt"

	dfapkg "github.com/timtadh/lexmachine/dfa"
	"github.com/timtadh/lexmachine/frontend"
	"github.com/timtadh/lexmachine/inst"
	"github.com/timtadh/lexmachine/machines"
	"github.com/timtadh/lexmachine/stream"
	"github.com/timtadh/lexmachine/stream_machines"
)

type pattern struct {
	regex  []byte
	action Action
}

// An Action is a function which get called when the Scanner finds a match
// during the lexing process. They turn a low level machines.Match struct into
// a token for the users program. As different compilers/interpretters/parsers
// have different needs Actions merely return an interface{}. This allows you
// to represent a token in anyway you wish. An example Token struct is provided
// above.
type Action func(scan Scanner, match *machines.Match) (interface{}, error)

// Lexer is a "builder" object which lets you construct a Scanner type which
// does the actual work of tokenizing (splitting up and categorizing) a byte
// string.  Get a new Lexer by calling the NewLexer() function. Add patterns to
// match (with their callbacks) by using the Add function. Finally, construct a
// scanner with Scanner to tokenizing a byte string.
type Lexer struct {
	patterns   []*pattern
	nfaMatches map[int]int // match_idx -> pat_idx
	dfaMatches map[int]int // match_idx -> pat_idx
	program    inst.Slice
	dfa        *dfapkg.DFA
}

// NewLexer constructs a new lexer object.
func NewLexer() *Lexer {
	return &Lexer{}
}

// TextScanner creates a scanner for a particular byte string from the lexer.
func (l *Lexer) TextScanner(text []byte) (*TextScanner, error) {
	if l.program == nil && l.dfa == nil {
		err := l.Compile()
		if err != nil {
			return nil, err
		}
	}

	// prevent the user from modifying the text under scan
	textCopy := make([]byte, len(text))
	copy(textCopy, text)

	var s *TextScanner
	if l.dfa != nil {
		s = &TextScanner{
			lexer:   l,
			matches: l.dfaMatches,
			scan:    machines.DFALexerEngine(l.dfa.Start, l.dfa.Error, l.dfa.Trans, l.dfa.Accepting, textCopy),
			Text:    textCopy,
			TC:      0,
		}
	} else {
		s = &TextScanner{
			lexer:   l,
			matches: l.nfaMatches,
			scan:    machines.LexerEngine(l.program, textCopy),
			Text:    textCopy,
			TC:      0,
		}
	}
	return s, nil
}

// StreamScanner creates a scanner for a particular stream from the lexer.
func (l *Lexer) StreamScanner(text stream.Stream) (*StreamScanner, error) {
	if l.program == nil && l.dfa == nil {
		err := l.Compile()
		if err != nil {
			return nil, err
		}
	}

	var s *StreamScanner
	if l.dfa != nil {
		s = &StreamScanner{
			lexer:   l,
			matches: l.dfaMatches,
			scan:    stream_machines.DFALexerEngine(l.dfa.Start, l.dfa.Error, l.dfa.Trans, l.dfa.Accepting, text),
			Text:    text,
		}
	} else {
		panic("not implemented")
	}
	return s, nil
}

// Add pattern to match on. When a match occurs during scanning the action
// function will be called by the Scanner to turn the low level machines.Match
// struct into a token.
func (l *Lexer) Add(regex []byte, action Action) {
	if l.program != nil {
		l.program = nil
	}
	l.patterns = append(l.patterns, &pattern{regex, action})
}

// Compile the supplied patterns to an DFA (default). You don't need to call
// this method (it is called automatically by Scanner). However, you may want to
// call this method if you construct a lexer once and then use it many times as
// it will precompile the lexing program.
func (l *Lexer) Compile() error {
	return l.CompileDFA()
}

func (l *Lexer) assembleAST() (frontend.AST, error) {
	asts := make([]frontend.AST, 0, len(l.patterns))
	for _, p := range l.patterns {
		ast, err := frontend.Parse(p.regex)
		if err != nil {
			return nil, err
		}
		asts = append(asts, ast)
	}
	lexast := asts[len(asts)-1]
	for i := len(asts) - 2; i >= 0; i-- {
		lexast = frontend.NewAltMatch(asts[i], lexast)
	}
	return lexast, nil
}

// CompileNFA compiles an NFA explicitly. If no DFA has been created (which is
// only created explicitly) this will be used by Scanners when they are created.
func (l *Lexer) CompileNFA() error {
	if len(l.patterns) == 0 {
		return fmt.Errorf("No patterns added")
	}
	if l.program != nil {
		return nil
	}
	lexast, err := l.assembleAST()
	if err != nil {
		return err
	}
	program, err := frontend.Generate(lexast)
	if err != nil {
		return err
	}

	l.program = program
	l.nfaMatches = make(map[int]int)

	ast := 0
	for i, instruction := range l.program {
		if instruction.Op == inst.MATCH {
			l.nfaMatches[i] = ast
			ast++
		}
	}

	if mes, err := l.matchesEmptyString(); err != nil {
		return err
	} else if mes {
		l.program = nil
		l.nfaMatches = nil
		return fmt.Errorf("One or more of the supplied patterns match the empty string")
	}

	return nil
}

// CompileDFA compiles an DFA explicitly. This will be used by Scanners when
// they are created.
func (l *Lexer) CompileDFA() error {
	if len(l.patterns) == 0 {
		return fmt.Errorf("No patterns added")
	}
	if l.dfa != nil {
		return nil
	}
	lexast, err := l.assembleAST()
	if err != nil {
		return err
	}
	dfa := dfapkg.Generate(lexast)
	l.dfa = dfa
	l.dfaMatches = make(map[int]int)
	for mid := range dfa.Matches {
		l.dfaMatches[mid] = mid
	}
	if mes, err := l.matchesEmptyString(); err != nil {
		return err
	} else if mes {
		l.dfa = nil
		l.dfaMatches = nil
		return fmt.Errorf("One or more of the supplied patterns match the empty string")
	}
	return nil
}

func (l *Lexer) matchesEmptyString() (bool, error) {
	s, err := l.TextScanner([]byte(""))
	if err != nil {
		return false, err
	}
	_, err, _ = s.Next()
	if ese, is := err.(*machines.EmptyMatchError); ese != nil && is {
		return true, nil
	}
	return false, nil
}
