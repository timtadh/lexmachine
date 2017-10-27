package lexmachine

import (
	"bytes"
	"fmt"
)

import (
	"github.com/timtadh/lexmachine/frontend"
	"github.com/timtadh/lexmachine/inst"
	"github.com/timtadh/lexmachine/machines"
)

// Token is an optional token representation you could use to represent the
// tokens produced by a lexer built with lexmachine.
//
// Here is an example for constructing a lexer Action which turns a
// machines.Match struct into a token using the scanners Token helper function.
//
//     func token(name string, tokenIds map[string]int) lex.Action {
//         return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
//             return s.Token(tokenIds[name], string(m.Bytes), m), nil
//         }
//     }
//
type Token struct {
	Type        int
	Value       interface{}
	Lexeme      []byte
	TC          int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

// Equals checks the equality of two tokens ignoring the Value field.
func (t *Token) Equals(other *Token) bool {
	if t == nil && other == nil {
		return true
	} else if t == nil {
		return false
	} else if other == nil {
		return false
	}
	return t.TC == other.TC &&
		t.StartLine == other.StartLine &&
		t.StartColumn == other.StartColumn &&
		t.EndLine == other.EndLine &&
		t.EndColumn == other.EndColumn &&
		bytes.Equal(t.Lexeme, other.Lexeme) &&
		t.Type == other.Type
}

// String formats the token in a human readable form.
func (t *Token) String() string {
	return fmt.Sprintf("%d %v %d (%d, %d)-(%d, %d)", t.Type, t.Value, t.TC, t.StartLine, t.StartColumn, t.EndLine, t.EndColumn)
}

// Actions are functions which get called when the Scanner finds a match during
// the lexing process. They turn a low level machines.Match struct into a token
// for the users program. As different compilers/interpretters/parsers have
// different needs Actions merely return an interface{}. This allows you to
// represent a token in anyway you wish. An example Token struct is provided
// above.
type Action func(scan *Scanner, match *machines.Match) (interface{}, error)

type pattern struct {
	regex  []byte
	action Action
}

// Lexer is a "builder" object which lets you construct a Scanner type which
// does the actual work of tokenizing (splitting up and categorizing) a byte
// string.  Get a new Lexer by calling the NewLexer() function. Add patterns to
// match (with their callbacks) by using the Add function. Finally, construct a
// scanner with Scanner to tokenizing a byte string.
type Lexer struct {
	patterns []*pattern
	matches  map[int]int "match_idx -> pat_idx"
	program  inst.InstSlice
}

// Scanner tokenizes a byte string based on the patterns provided to the lexer
// object which constructed the scanner. This object works as functional
// iterator using the Next method.
//
// Example
//
//     lexer, err := CreateLexer()
//     if err != nil {
//         return err
//     }
//     scanner, err := lexer.Scanner(someBytes)
//     if err != nil {
//         return err
//     }
//     for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
//         if err != nil {
//             return err
//         }
//         fmt.Println(tok)
//     }
//
type Scanner struct {
	lexer    *Lexer
	scan     machines.Scanner
	Text     []byte
	TC       int
	pTC      int
	s_line   int
	s_column int
	e_line   int
	e_column int
}

// Next iterates through the string being scanned returning one token at a time
// until either an error is encountered or the end of the string is reached.
// The token is returned by the tok value. An error is indicated by err.
// Finally, eos (a bool) indicates the End Of String when it returns as true.
//
// Example
//
//     for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
//         if err != nil {
//             // handle the error and exit the loop. For example:
//             return err
//         }
//         // do some processing on tok or store it somewhere. eg.
//         fmt.Println(tok)
//     }
//
func (s *Scanner) Next() (tok interface{}, err error, eos bool) {
	var token interface{} = nil
	for token == nil {
		tc, match, err, scan := s.scan(s.TC)
		if scan == nil {
			return nil, nil, true
		} else if err != nil {
			return nil, err, false
		} else if match == nil {
			return nil, fmt.Errorf("No match but no error"), false
		}
		s.scan = scan
		s.pTC = s.TC
		s.TC = tc
		s.s_line = match.StartLine
		s.s_column = match.StartColumn
		s.e_line = match.EndLine
		s.e_column = match.EndColumn

		pattern := s.lexer.patterns[s.lexer.matches[match.PC]]
		token, err = pattern.action(s, match)
		if err != nil {
			return nil, err, false
		}
	}
	return token, nil, false
}

// Token is a helper function for constructing a Token type inside of a Action.
func (s *Scanner) Token(typ int, value interface{}, m *machines.Match) *Token {
	return &Token{
		Type:        typ,
		Value:       value,
		Lexeme:      m.Bytes,
		TC:          m.TC,
		StartLine:   m.StartLine,
		StartColumn: m.StartColumn,
		EndLine:     m.EndLine,
		EndColumn:   m.EndColumn,
	}
}

// NewLexer constructs a new lexer object.
func NewLexer() *Lexer {
	return &Lexer{}
}

// Scanner creates a scanner for a particular byte string from the lexer.
func (l *Lexer) Scanner(text []byte) (*Scanner, error) {
	err := l.Compile()
	if err != nil {
		return nil, err
	}

	// prevent the user from modifying the text under scan
	text_copy := make([]byte, len(text))
	copy(text_copy, text)

	return &Scanner{
		lexer: l,
		scan:  machines.LexerEngine(l.program, text_copy),
		Text:  text_copy,
		TC:    0,
	}, nil
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

// Compiles the supplied patterns. You don't need to call this method (it is
// called automatically by Scanner). However, you may want to call this method
// if you construct a lexer once and then use it many times as it will
// precompile the lexing program.
func (l *Lexer) Compile() error {
	if len(l.patterns) == 0 {
		return fmt.Errorf("No patterns added")
	}
	if l.program != nil {
		return nil
	}

	asts := make([]frontend.AST, 0, len(l.patterns))
	for _, p := range l.patterns {
		ast, err := frontend.Parse(p.regex)
		if err != nil {
			return err
		}
		asts = append(asts, ast)
	}

	lexast := asts[len(asts)-1]
	for i := len(asts) - 2; i >= 0; i-- {
		lexast = frontend.NewAltMatch(asts[i], lexast)
	}

	program, err := frontend.Generate(lexast)
	if err != nil {
		return err
	}

	l.program = program
	l.matches = make(map[int]int)

	ast := 0
	for i, instruction := range l.program {
		if instruction.Op == inst.MATCH {
			l.matches[i] = ast
			ast += 1
		}
	}

	if len(asts) != ast {
		panic("len(asts) != ast")
	}

	return nil
}
