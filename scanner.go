package lexmachine

import (
	"fmt"

	"github.com/timtadh/lexmachine/machines"
)

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
	lexer   *Lexer
	matches map[int]int
	scan    machines.Scanner
	Text    []byte
	TC      int
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
// One useful error type which could be returned by Next() is a
// match.UnconsumedInput which provides the position information for where in
// the text the scanning failed.
//
// For more information on functional iterators see:
// http://hackthology.com/functional-iteration-in-go.html
func (s *Scanner) Next() (tok interface{}, err error, eos bool) {
	var token interface{}
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
		s.TC = tc

		pattern := s.lexer.patterns[s.matches[match.PC]]
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
