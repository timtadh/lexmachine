package lexmachine

import "testing"

import (
	"strconv"
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
		func(scan *Scanner, match []byte)(*Token, error) {
			return &Token{PRINT, nil, match}, nil
		},
	)
	lexer.Add(
		[]byte("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*"),
		func(scan *Scanner, match []byte)(*Token, error) {
			return &Token{NAME, string(match), match}, nil
		},
	)
	lexer.Add(
		[]byte("="),
		func(scan *Scanner, match []byte)(*Token, error) {
			return &Token{EQUALS, nil, match}, nil
		},
	)
	lexer.Add(
		[]byte("[0-9]+"),
		func(scan *Scanner, match []byte)(*Token, error) {
			i, err := strconv.Atoi(string(match))
			if err != nil {
				return nil, err
			}
			return &Token{NUMBER, i, match}, nil
		},
	)
	lexer.Add(
		[]byte("( |\t|\n)"),
		func(scan *Scanner, match []byte)(*Token, error) {
			// skip white space
			return nil, nil
		},
	)

	scanner, err := lexer.Scanner([]byte(`
		name = 10
		print name
		print fred
		name =12
		printname = 13
		print printname
	`))
	if err != nil {
		t.Error(err)
	}

	expected := []int{
		NAME, EQUALS, NUMBER,
		PRINT, NAME,
		PRINT, NAME,
		NAME, EQUALS, NUMBER,
		NAME, EQUALS, NUMBER,
		PRINT, NAME,
	}

	t.Log(lexer.program)

	i := 0
	for tk, err, eof := scanner.Scan(); !eof; tk, err, eof = scanner.Scan() {
		tok := tk.(*Token)
		t.Log(tok)
		if err != nil {
			t.Fatal(err)
		}
		if tok.Type != expected[i] {
			t.Errorf("got wrong token %d != %d", tok.Type, expected[i])
		}
		i += 1
	}
}

