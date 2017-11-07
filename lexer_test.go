package lexmachine

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/timtadh/data-structures/test"
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
		func(scan *Scanner, match *machines.Match) (interface{}, error) {
			return scan.Token(PRINT, nil, match), nil
		},
	)
	lexer.Add(
		[]byte("([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*"),
		func(scan *Scanner, match *machines.Match) (interface{}, error) {
			return scan.Token(NAME, string(match.Bytes), match), nil
		},
	)
	lexer.Add(
		[]byte("="),
		func(scan *Scanner, match *machines.Match) (interface{}, error) {
			return scan.Token(EQUALS, nil, match), nil
		},
	)
	lexer.Add(
		[]byte("[0-9]+"),
		func(scan *Scanner, match *machines.Match) (interface{}, error) {
			i, err := strconv.Atoi(string(match.Bytes))
			if err != nil {
				return nil, err
			}
			return scan.Token(NUMBER, i, match), nil
		},
	)
	lexer.Add(
		[]byte("( |\t|\n)"),
		func(scan *Scanner, match *machines.Match) (interface{}, error) {
			// skip white space
			return nil, nil
		},
	)
	lexer.Add(
		[]byte("//[^\n]*\n"),
		func(scan *Scanner, match *machines.Match) (interface{}, error) {
			// skip white space
			return nil, nil
		},
	)
	lexer.Add(
		[]byte("/\\*"),
		func(scan *Scanner, match *machines.Match) (interface{}, error) {
			for tc := scan.TC; tc < len(scan.Text); tc++ {
				if scan.Text[tc] == '\\' {
					// the next character is skipped
					tc++
				} else if scan.Text[tc] == '*' && tc+1 < len(scan.Text) {
					if scan.Text[tc+1] == '/' {
						scan.TC = tc + 2
						return nil, nil
					}
				}
			}
			return nil,
				fmt.Errorf("unclosed comment starting at %d, (%d, %d)",
					match.TC, match.StartLine, match.StartColumn)
		},
	)

	text := []byte(`
		name = 10
		print name
		print fred
		name =12
		// asdf comment
		/*awef  oiwe
		 ooiwje \*/ weoi
		 weoi*/ printname = 13
		print printname
	`)

	expected := []*Token{
		{NAME, "name", []byte("name"), 3, 2, 3, 2, 6},
		{EQUALS, nil, []byte("="), 8, 2, 8, 2, 8},
		{NUMBER, 10, []byte("10"), 10, 2, 10, 2, 11},
		{PRINT, nil, []byte("print"), 15, 3, 3, 3, 7},
		{NAME, "name", []byte("name"), 21, 3, 9, 3, 12},
		{PRINT, nil, []byte("print"), 28, 4, 3, 4, 7},
		{NAME, "fred", []byte("fred"), 34, 4, 9, 4, 12},
		{NAME, "name", []byte("name"), 41, 5, 3, 5, 6},
		{EQUALS, nil, []byte("="), 46, 5, 8, 5, 8},
		{NUMBER, 12, []byte("12"), 47, 5, 9, 5, 10},
		{NAME, "printname", []byte("printname"), 112, 9, 11, 9, 19},
		{EQUALS, nil, []byte("="), 122, 9, 21, 9, 21},
		{NUMBER, 13, []byte("13"), 124, 9, 23, 9, 24},
		{PRINT, nil, []byte("print"), 129, 10, 3, 10, 7},
		{NAME, "printname", []byte("printname"), 135, 10, 9, 10, 17},
	}

	// first do the test with the NFA
	err := lexer.CompileNFA()
	if err != nil {
		t.Error(err)
	}

	scanner, err := lexer.Scanner(text)
	if err != nil {
		t.Error(err)
		t.Log(lexer.program.Serialize())
	}

	i := 0
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if err != nil {
			t.Fatal(err)
		}
		tok := tk.(*Token)
		if !tok.Equals(expected[i]) {
			t.Errorf("got wrong token got %v, expected %v", tok, expected[i])
		}
		i++
	}

	lexer.program = nil
	lexer.nfaMatches = nil

	// first do the test with the DFA
	err = lexer.CompileDFA()
	if err != nil {
		t.Error(err)
	}

	scanner, err = lexer.Scanner(text)
	if err != nil {
		t.Error(err)
		t.Log(lexer.program.Serialize())
	}

	i = 0
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if err != nil {
			t.Fatal(err)
		}
		tok := tk.(*Token)
		if !tok.Equals(expected[i]) {
			t.Errorf("got wrong token got %v, expected %v", tok, expected[i])
		}
		i++
	}
}

func TestPartialLexer(x *testing.T) {
	t := (*test.T)(x)
	text := `
	require 'config.php';
	class Jikan
	{
		public $response = [];
		public function __construct() {
			return $this;
		}
		/*
		 * Anime
		 */
		public function Anime(String $id = null, Array $extend = []) {
			$this->response = (array) (new Get\Anime($id, $extend))->response;
			return $this;
		}
		/*
		 * Manga
		 */
		public function Manga(String $id = null, Array $extend = []) {
			$this->response = (array) (new Get\Manga($id, $extend))->response;
			return $this;
		}`
	tokens := []string{
		"ERROR",
		"INCLUDE",
		"COMMENT",
		"STRING",
		"IDENT",
		"FUNC",
		"CLASS",
		"OP",
	}
	tokmap := make(map[string]int)
	for id, name := range tokens {
		tokmap[name] = id
	}
	expected := []int{
		tokmap["INCLUDE"],
		tokmap["STRING"],
		tokmap["CLASS"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["FUNC"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["COMMENT"],
		tokmap["IDENT"],
		tokmap["FUNC"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["COMMENT"],
		tokmap["IDENT"],
		tokmap["FUNC"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["STRING"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["COMMENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["COMMENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["OP"],
		tokmap["IDENT"],
		tokmap["IDENT"],
		tokmap["IDENT"],
	}

	getToken := func(tokenType int) Action {
		return func(s *Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
	}
	var lexer *Lexer = NewLexer()
	lexer.Add([]byte("import|require"), getToken(tokmap["INCLUDE"]))
	lexer.Add([]byte("function"), getToken(tokmap["FUNC"]))
	lexer.Add([]byte("class"), getToken(tokmap["CLASS"]))
	lexer.Add([]byte("\"[^\\\"]*\"|'[^']*'|`[^`]*`"), getToken(tokmap["STRING"]))
	lexer.Add([]byte("//[^\n]*\n?|/\\*([^*]|\r|\n|(\\*+([^*/]|\r|\n)))*\\*+/"), getToken(tokmap["COMMENT"]))
	lexer.Add([]byte("[A-Za-z$][A-Za-z0-9$]+"), getToken(tokmap["IDENT"]))
	lexer.Add([]byte(">=|<=|=|>|<|\\|\\||&&"), getToken(tokmap["OP"]))
	scan := func(lexer *Lexer) {
		scanner, err := lexer.Scanner([]byte(text))
		t.AssertNil(err)
		i := 0
		for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
			if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
				scanner.TC = ui.FailTC
				// t.Log(ui)
			} else if err != nil {
				t.Fatal(err)
			} else {
				t.Assert(tk.(*Token).Type == expected[i],
					"expected %v got %v: %v", tokens[expected[i]], tokens[tk.(*Token).Type], tk)
				i++
			}
		}
	}
	t.AssertNil(lexer.CompileNFA())
	scan(lexer)
	lexer.program = nil
	lexer.nfaMatches = nil
	t.AssertNil(lexer.CompileDFA())
	scan(lexer)
}
