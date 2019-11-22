package lexmachine

import (
	"bytes"
	"fmt"
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
	return fmt.Sprintf("%d %q %d (%d, %d)-(%d, %d)", t.Type, t.Value, t.TC, t.StartLine, t.StartColumn, t.EndLine, t.EndColumn)
}
