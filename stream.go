package lexmachine

import (
	"fmt"

	"github.com/timtadh/lexmachine/machines"
	"github.com/timtadh/lexmachine/stream"
	"github.com/timtadh/lexmachine/stream_machines"
)

// StreamScanner tokenizes a stream of bytes (see stream.Stream) which can be
// constructed from an io.Reader. This object work analogously to the regular
// Scanner. Note: if the stream you are scanning fits in memory using the
// regular Scanner is likely more efficient. Finally, stream.Stream objects can
// only advance the text forwards so an Action cannot move the text counter
// backwards (as is possible with Scanner).
type StreamScanner struct {
	lexer   *Lexer
	matches map[int]int
	scan    stream_machines.Scanner
	Text    stream.Stream
	buf     *StreamBuffer
}

func (s *StreamScanner) Buffer() Buffer {
	if s.buf == nil {
		panic(fmt.Errorf("Buffer called outside of an Action"))
	}
	return s.buf
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
func (s *StreamScanner) Next() (tok interface{}, err error, eos bool) {
	var token interface{}
	for token == nil {
		match, err, scan := s.scan()
		if scan == nil {
			return nil, nil, true
		} else if err != nil {
			return nil, err, false
		} else if match == nil {
			return nil, fmt.Errorf("No match but no error"), false
		}
		s.scan = scan

		s.buf = streamBuffer(s.Text)
		pattern := s.lexer.patterns[s.matches[match.PC]]
		token, err = pattern.action(s, match)
		s.buf.finalize()
		s.buf = nil
		if err != nil {
			return nil, err, false
		}
	}
	return token, nil, false
}

// Token is a helper function for constructing a Token type inside of a Action.
func (s *StreamScanner) Token(typ int, value interface{}, m *machines.Match) *Token {
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
