package lexmachine

import "testing"

import (
	"fmt"
	"strconv"
)

import (
	"github.com/timtadh/lexmachine/machines"
)

func TestSimple(t *testing.T) {
	const (
		NAME = iota
		EQUALS
		NUMBER
		PRINT
	)
	lexer := NewLexer()

	lexer.Add(
		[]byte("print"),
		func(scan *Scanner, match *machines.Match)(interface{}, error) {
			return scan.Token(PRINT, nil, match.Bytes), nil
		},
	)
	lexer.Add(
		[]byte("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*"),
		func(scan *Scanner, match *machines.Match)(interface{}, error) {
			return scan.Token(NAME, string(match.Bytes), match.Bytes), nil
		},
	)
	lexer.Add(
		[]byte("="),
		func(scan *Scanner, match *machines.Match)(interface{}, error) {
			return scan.Token(EQUALS, nil, match.Bytes), nil
		},
	)
	lexer.Add(
		[]byte("[0-9]+"),
		func(scan *Scanner, match *machines.Match)(interface{}, error) {
			i, err := strconv.Atoi(string(match.Bytes))
			if err != nil {
				return nil, err
			}
			return scan.Token(NUMBER, i, match.Bytes), nil
		},
	)
	lexer.Add(
		[]byte("( |\t|\n)"),
		func(scan *Scanner, match *machines.Match)(interface{}, error) {
			// skip white space
			return nil, nil
		},
	)
	lexer.Add(
		[]byte("//[^\n]*\n"),
		func(scan *Scanner, match *machines.Match)(interface{}, error) {
			// skip white space
			return nil, nil
		},
	)
	lexer.Add(
		[]byte("/\\*"),
		func(scan *Scanner, match *machines.Match)(interface{}, error) {
			for tc := scan.TC; tc < len(scan.Text); tc++ {
				if scan.Text[tc] == '\\' {
					// the next character is skipped
					tc++
				} else if scan.Text[tc] == '*' && tc+1 < len(scan.Text) {
					if scan.Text[tc+1] == '/' {
						scan.TC = tc+2
						return nil, nil
					}
				}
			}
			return nil,
				fmt.Errorf("unclosed comment starting at %d, (%d, %d)",
					match.TC, match.Line, match.Column)
		},
	)

	scanner, err := lexer.Scanner([]byte(`
		name = 10
		print name
		print fred
		name =12
		// asdf comment
		/*awef  oiwe
		 ooiwje \*/ weoi
		 weoi*/ printname = 13
		print printname
	`))
	if err != nil {
		t.Error(err)
	}

	expected := []*Token{
		&Token{NAME, "name", []byte("name"), 3, 2, 3},
			&Token{EQUALS, nil, []byte("="), 8, 2, 8},
			&Token{NUMBER, 10, []byte("10"), 10, 2, 10},
		&Token{PRINT, nil, []byte("print"), 15, 3, 3},
			&Token{NAME, "name", []byte("name"), 21, 3, 9},
		&Token{PRINT, nil, []byte("print"), 28, 4, 3},
			&Token{NAME, "fred", []byte("fred"), 34, 4, 9},
		&Token{NAME, "name", []byte("name"), 41, 5, 3},
			&Token{EQUALS, nil, []byte("="), 46, 5, 8},
			&Token{NUMBER, 12, []byte("12"), 57, 5, 9},
		&Token{NAME, "printname", []byte("printname"), 94, 8, 11},
			&Token{EQUALS, nil, []byte("="), 104, 8, 21},
			&Token{NUMBER, 13, []byte("13"), 106, 8, 23},
		&Token{PRINT, nil, []byte("print"), 111, 9, 3},
			&Token{NAME, "printname", []byte("printname"), 117, 9, 9},
	}

	t.Log(lexer.program)

	i := 0
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if err != nil {
			t.Fatal(err)
		}
		tok := tk.(*Token)
		t.Log(tok)
		if tok.Equals(expected[i]) {
			t.Errorf("got wrong token got %v, expected %v", tok, expected[i])
		}
		i += 1
	}
}

