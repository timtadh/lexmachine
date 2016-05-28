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

func (self *Token) Equals(other *Token) bool {
	if self == nil && other == nil {
		return true
	} else if self == nil {
		return false
	} else if other == nil {
		return false
	}
	return self.TC == other.TC &&
		self.StartLine == other.StartLine &&
		self.StartColumn == other.StartColumn &&
		self.EndLine == other.EndLine &&
		self.EndColumn == other.EndColumn &&
		bytes.Equal(self.Lexeme, other.Lexeme) &&
		self.Type == other.Type
}

func (self *Token) String() string {
	return fmt.Sprintf("%d %v %d (%d, %d)-(%d, %d)", self.Type, self.Value, self.TC, self.StartLine, self.StartColumn, self.EndLine, self.EndColumn)
}

type Action func(scan *Scanner, match *machines.Match) (interface{}, error)

type Pattern struct {
	regex  []byte
	action Action
}

type Lexer struct {
	patterns []*Pattern
	matches  map[int]int "match_idx -> pat_idx"
	program  inst.InstSlice
}

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

func (self *Scanner) Next() (tok interface{}, err error, eof bool) {
	var token interface{} = nil
	for token == nil {
		tc, match, err, scan := self.scan(self.TC)
		if scan == nil {
			return nil, nil, true
		} else if err != nil {
			return nil, err, false
		} else if match == nil {
			return nil, fmt.Errorf("No match but no error"), false
		}
		self.scan = scan
		self.pTC = self.TC
		self.TC = tc
		self.s_line = match.StartLine
		self.s_column = match.StartColumn
		self.e_line = match.EndLine
		self.e_column = match.EndColumn

		pattern := self.lexer.patterns[self.lexer.matches[match.PC]]
		token, err = pattern.action(self, match)
		if err != nil {
			return nil, err, false
		}
	}
	return token, nil, false
}

func (self *Scanner) Token(typ int, value interface{}, m *machines.Match) *Token {
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

func NewLexer() *Lexer {
	return &Lexer{}
}

func (self *Lexer) Scanner(text []byte) (*Scanner, error) {
	err := self.Compile()
	if err != nil {
		return nil, err
	}

	scan := machines.LexerEngine(self.program, text)

	// prevent the user from modifying the text under scan
	text_copy := make([]byte, len(text))
	copy(text_copy, text)

	return &Scanner{
		lexer: self,
		scan:  scan,
		Text:  text_copy,
		TC:    0,
	}, nil
}

func (self *Lexer) Add(regex []byte, action Action) {
	self.patterns = append(self.patterns, &Pattern{regex, action})
}

// Compiles the supplied patterns. You don't need
func (self *Lexer) Compile() error {
	if len(self.patterns) == 0 {
		return fmt.Errorf("No patterns added")
	}
	if self.program != nil {
		return nil
	}

	asts := make([]frontend.AST, 0, len(self.patterns))
	for _, p := range self.patterns {
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

	self.program = program
	self.matches = make(map[int]int)

	ast := 0
	for i, instruction := range self.program {
		if instruction.Op == inst.MATCH {
			self.matches[i] = ast
			ast += 1
		}
	}

	if len(asts) != ast {
		panic("len(asts) != ast")
	}

	return nil
}
