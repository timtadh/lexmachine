//go:generate sh -c "if [[ -z \"`which goyacc`\" ]]; then go get -u golang.org/x/tools/cmd/goyacc || exit 1; fi; if [[ -f y.go ]]; then rm y.*; fi; goyacc sensors.y"
//
// Package golex implements the same lexer as examples/sensors. However, it
// shows how to conform to the goyacc's expected interface:
//
// type yyLexer interface {
//    Lex(lval *yySymType) (tokenType int)
//    Error(message string)
// }
//
// You define yySymType. The yyLexer type is defined by the generated code
// from goyacc. The tokenType is the token identifier. The expectation is
// the token id's are shared between what is defined in this package and
// the parser definition in parser.y.
//
// To generate the parser (and make this all work) run:
//
// go generate github.com/timtadh/lexmachine/examples/golex
//
package main

import (
	"fmt"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

type golex struct {
	*lexmachine.Scanner
	stmts []*Node
}

// Construct a new golex from a lexer object and the text to parse.
func newGoLex(lexer *lexmachine.Lexer, text []byte) (*golex, error) {
	scan, err := lexer.Scanner(text)
	if err != nil {
		return nil, err
	}
	return &golex{Scanner: scan}, nil
}

// Lex implements yyLexer's interface for getting the next token. It returns the
// token type as an integer. The tokens should be defined in the $parser.y file.
// The actual number returned will be >= yyPrivate - 1 which is the range for
// custom token names.
func (g *golex) Lex(lval *yySymType) (tokenType int) {
	s := g.Scanner
	tok, err, eof := s.Next()
	if err != nil {
		g.Error(err.Error())
	} else if eof {
		return -1 // signals EOF to goyacc's yyParse
	}
	lval.token = tok.(*lexmachine.Token)
	// To return the correct number for goyacc you must add yyPrivate - 1 to
	// put the value into the correct range.
	return lval.token.Type + yyPrivate - 1
}

// Error implements the error handling for if there is a parse error of any
// kind. This implementation panics. There may be no better way to hand errors
// from goyacc. I recommend you use defer ... recover() to handle this where
// you call into the parser.
func (l *golex) Error(message string) {
	// is there a better way to handle this in the context of goyacc?
	panic(fmt.Errorf(message))
}

// newLexer constructs the lexer for you. Only call this once.
func newLexer() *lexmachine.Lexer {
	// build the token map from yyToknames produced by goyacc from the %token
	// directives.
	tokmap := make(map[string]int)
	for id, name := range yyToknames {
		tokmap[name] = id
	}
	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
	}
	skip := func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return nil, nil
	}
	var lexer = lexmachine.NewLexer()
	lexer.Add([]byte("@"), getToken(tokmap["AT"]))
	lexer.Add([]byte(`\+`), getToken(tokmap["PLUS"]))
	lexer.Add([]byte(`\*`), getToken(tokmap["STAR"]))
	lexer.Add([]byte("-"), getToken(tokmap["DASH"]))
	lexer.Add([]byte("/"), getToken(tokmap["SLASH"]))
	lexer.Add([]byte("\\"), getToken(tokmap["BACKSLASH"]))
	lexer.Add([]byte(`\^`), getToken(tokmap["CARROT"]))
	lexer.Add([]byte("`"), getToken(tokmap["BACKTICK"]))
	lexer.Add([]byte(","), getToken(tokmap["COMMA"]))
	lexer.Add([]byte(`\(`), getToken(tokmap["LPAREN"]))
	lexer.Add([]byte(`\)`), getToken(tokmap["RPAREN"]))
	lexer.Add([]byte("bus"), getToken(tokmap["BUS"]))
	lexer.Add([]byte("chip"), getToken(tokmap["CHIP"]))
	lexer.Add([]byte("label"), getToken(tokmap["LABEL"]))
	lexer.Add([]byte("compute"), getToken(tokmap["COMPUTE"]))
	lexer.Add([]byte("ignore"), getToken(tokmap["IGNORE"]))
	lexer.Add([]byte("set"), getToken(tokmap["SET"]))
	lexer.Add([]byte(`[0-9]*\.?[0-9]+`), getToken(tokmap["NUMBER"]))
	lexer.Add([]byte(`[a-zA-Z_][a-zA-Z0-9_]*`), getToken(tokmap["NAME"]))
	lexer.Add([]byte(`"[^"]*"`), getToken(tokmap["NAME"]))
	lexer.Add([]byte(`\\\n`), skip) // skip backslash newline
	lexer.Add([]byte(`\n`), getToken(tokmap["NEWLINE"]))
	lexer.Add([]byte(`#[^\n]*`), skip)
	lexer.Add([]byte(` |\t`), skip)
	err := lexer.Compile()
	if err != nil {
		panic(err)
	}
	return lexer
}
