package lexmachine

import (
	"fmt"
	"reflect"
	"bytes"
)

import (
	"github.com/timtadh/lexmachine/frontend"
	"github.com/timtadh/lexmachine/machines"
	"github.com/timtadh/lexmachine/inst"
)

type Token struct {
	Type int
	Value interface{}
	Lexeme []byte
	TC int
	Line int
	Column int
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
			self.Line == other.Line &&
			self.Column == other.Column &&
			bytes.Equal(self.Lexeme, other.Lexeme) &&
			self.Type == other.Type &&
			reflect.DeepEqual(self.Value, other.Value)
}

func (self *Token) String() string {
	return fmt.Sprintf("%d %v (%d, %d)", self.Type, self.Value, self.Line, self.Column)
}

type Action func(scan *Scanner, match *machines.Match) (*Token, error)

type Pattern struct {
	regex []byte
	action Action
}

type Lexer struct {
	patterns []*Pattern
	matches map[int]int "match_idx -> pat_idx"
	program inst.InstSlice
}

type Scanner struct {
	lexer *Lexer
	scan machines.Scanner
	Text []byte
	TC int
	line int
	column int
}

func (self *Scanner) Next() (tok interface{}, err error, eof bool) {
	tc, match, err, scan := self.scan(self.TC)
	if scan == nil {
		return nil, nil, true
	} else if err != nil {
		return nil, err, false
	} else if match == nil {
		return nil, fmt.Errorf("No match but no error"), false
	}
	self.scan = scan
	self.TC = tc
	self.line = match.Line
	self.column = match.Column

	pattern := self.lexer.patterns[self.lexer.matches[match.PC]]
	token, err := pattern.action(self, match)
	if err != nil {
		return nil, err, false
	} else if token == nil {
		return self.Next()
	}

	return token, nil, false
}

func (self *Scanner) Token(typ int, value interface{}, lexeme []byte) *Token {
	return &Token{
		Type: typ,
		Value: value,
		Lexeme: lexeme,
		TC: self.TC,
		Line: self.line,
		Column: self.column,
	}
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (self *Lexer) Scanner(text []byte) (*Scanner, error) {
	if self.program == nil || len(self.patterns) != len(self.matches) {
		err := self.Compile()
		if err != nil {
			return nil, err
		}
	}

	scan := machines.LexerEngine(self.program, text)

	// prevent the user from modifying the text under scan
	text_copy := make([]byte, len(text))
	copy(text_copy, text)

	return &Scanner{
		lexer: self,
		scan: scan,
		Text: text_copy,
		TC: 0,
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

	asts := make([]frontend.AST, 0, len(self.patterns))
	for _, p := range self.patterns {
		ast, err := frontend.Parse(p.regex)
		if err != nil {
			return err
		}
		asts = append(asts, ast)
	}

	lexast := asts[len(asts)-1]
	for i := len(asts)-2; i >= 0; i-- {
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

