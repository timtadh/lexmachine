package lexmachine

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/timtadh/data-structures/test"
	"github.com/timtadh/lexmachine/machines"
)

func TestSimple(x *testing.T) {
	t := (*test.T)(x)
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

	scan := func(lexer *Lexer) {
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
	}

	// first do the test with the NFA
	t.AssertNil(lexer.CompileNFA())
	scan(lexer)

	// then do the test with the DFA
	lexer.program = nil
	lexer.nfaMatches = nil
	t.AssertNil(lexer.CompileDFA())
	scan(lexer)
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
		tokmap["INCLUDE"], tokmap["STRING"], tokmap["CLASS"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"],
		tokmap["FUNC"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["COMMENT"], tokmap["IDENT"], tokmap["FUNC"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"],
		tokmap["OP"], tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["COMMENT"], tokmap["IDENT"], tokmap["FUNC"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["OP"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["OP"],
		tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"], tokmap["OP"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["STRING"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["OP"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["COMMENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["OP"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["OP"], tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"],
		tokmap["OP"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["OP"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["COMMENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"], tokmap["OP"],
		tokmap["IDENT"], tokmap["OP"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
		tokmap["OP"], tokmap["IDENT"], tokmap["IDENT"], tokmap["IDENT"],
	}

	getToken := func(tokenType int) Action {
		return func(s *Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
	}
	var lexer = NewLexer()
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
				t.Log(ui)
			} else if err != nil {
				t.Fatal(err)
			} else {
				t.Logf("%v: %v", tokens[tk.(*Token).Type], tk)
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

func TestRegression(t *testing.T) {
	token := func(name string) Action {
		return func(s *Scanner, m *machines.Match) (interface{}, error) {
			return fmt.Sprintf("%v:%q", name, string(m.Bytes)), nil
		}
	}

	newLexer := func() *Lexer {
		lexer := NewLexer()
		lexer.Add([]byte("true"), token("TRUE"))
		lexer.Add([]byte("( |\t|\n|\r)+"), token("SPACE"))
		return lexer
	}

	tests := []struct {
		text   string
		tokens int
	}{
		{`true`, 1},
		{`true `, 2},
	}

	runTest := func(lexer *Lexer) {
		for _, test := range tests {
			scanner, err := lexer.Scanner([]byte(test.text))
			if err != nil {
				t.Fatal(err)
			}

			found := 0
			tok, err, eos := scanner.Next()
			for ; !eos; tok, err, eos = scanner.Next() {
				if err != nil {
					t.Fatal(err)
				}
				fmt.Printf("Token: %v\n", tok)
				found++
			}
			if found != test.tokens {
				t.Errorf("Expected exactly %v tokens got %v, ===\nErr: %v\nEOS: %v\nTC: %d\n", test.tokens, found, err, eos, scanner.TC)
			}
		}
	}
	{
		lexer := newLexer()
		if err := lexer.CompileNFA(); err != nil {
			t.Fatal(err)
		}
		runTest(lexer)
	}
	{
		lexer := newLexer()
		if err := lexer.CompileDFA(); err != nil {
			t.Fatal(err)
		}
		runTest(lexer)
	}
}

func TestRegression2(t *testing.T) {

	text := `# dhcpd.conf
#
# Sample configuration file for ISC dhcpd
#

# option definitions common to all supported networks...
option domain-name "example.org";
option domain-name-servers ns1.example.org, ns2.example.org;

default-lease-time 600;
max-lease-time 7200;

# The ddns-updates-style parameter controls whether or not the server will
# attempt to do a DNS update when a lease is confirmed. We default to the
# behavior of the version 2 packages ('none', since DHCP v2 didn't
# have support for DDNS.)
ddns-update-style none;

# If this DHCP server is the official DHCP server for the local
# network, the authoritative directive should be uncommented.
#authoritative;
`

	literals := []string{
		"{",
		"}",
		";",
		",",
	}
	tokens := []string{
		"COMMENT",
		"ID",
	}
	tokens = append(tokens, literals...)
	tokenIds := map[string]int{}
	for i, tok := range tokens {
		tokenIds[tok] = i
	}
	newLexer := func() *Lexer {
		lex := NewLexer()

		skip := func(*Scanner, *machines.Match) (interface{}, error) {
			return nil, nil
		}
		token := func(name string) Action {
			return func(s *Scanner, m *machines.Match) (interface{}, error) {
				return s.Token(tokenIds[name], string(m.Bytes), m), nil
			}
		}

		lex.Add([]byte(`#[^\n]*\n?`), token("COMMENT"))
		lex.Add([]byte(`([a-z]|[A-Z]|[0-9]|_|\-|\.)+`), token("ID"))
		lex.Add([]byte(`"([^\\"]|(\\.))*"`), token("ID"))
		lex.Add([]byte("[\n \t]"), skip)
		for _, lit := range literals {
			lex.Add([]byte(lit), token(lit))
		}
		return lex
	}

	runTest := func(lexer *Lexer) {
		scanner, err := lexer.Scanner([]byte(text))
		if err != nil {
			return
		}
		for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {
			if err != nil {
				t.Fatal(err)
				break
			}
			token := tok.(*Token)
			fmt.Printf("%-7v | %-10v | %v:%v-%v:%v\n",
				tokens[token.Type],
				strings.TrimSpace(string(token.Lexeme)),
				token.StartLine,
				token.StartColumn,
				token.EndLine,
				token.EndColumn)
		}
	}
	{
		lexer := newLexer()
		if err := lexer.CompileNFA(); err != nil {
			t.Fatal(err)
		}
		runTest(lexer)
	}
	{
		lexer := newLexer()
		if err := lexer.CompileDFA(); err != nil {
			t.Fatal(err)
		}
		runTest(lexer)
	}
}

func TestPythonStrings(t *testing.T) {
	tokens := []string{
		"UNDEF",
		"TRUE",
		"SINGLE_STRING",
		"TRIPLE_STRING",
		"TRIPLE_STRING2",
		"TY_STRING",
		"SPACE",
	}
	tokenIds := map[string]int{}
	for i, tok := range tokens {
		tokenIds[tok] = i
	}
	skip := func(*Scanner, *machines.Match) (interface{}, error) {
		return nil, nil
	}
	token := func(name string) Action {
		return func(s *Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenIds[name], string(m.Bytes), m), nil
		}
	}

	newLexer := func() *Lexer {
		lexer := NewLexer()
		lexer.Add([]byte("true"), token("TRUE"))
		lexer.Add([]byte(`'''([^\\']|(\\.))*'''`), token("TRIPLE_STRING"))
		lexer.Add([]byte(`"""([^\\"]|(\\.))*"""`), token("TRIPLE_STRING"))
		lexer.Add([]byte(`"([^\\"]|(\\.))*"`), token("SINGLE_STRING"))
		lexer.Add([]byte(`'([^\\']|(\\.))*'`), token("SINGLE_STRING"))
		lexer.Add([]byte("( |\t|\n|\r)+"), skip)
		return lexer
	}

	tests := []struct {
		text   string
		tokens int
	}{
		{`'''hi'''`, 1},
		{`"""hi"""`, 1},
		{`"hi"`, 1},
		{`'hi'`, 1},
		{`''`, 1},
		{`""`, 1},
		{`"""  .  .
			hello
		"""`, 1},
		{`'''' ''''`, 4},
		{`''''''`, 1},
		{`""""""`, 1},
		{`"""""" """
		hi there""" "wizard" true`, 4},
	}

	runTest := func(lexer *Lexer) {
		for _, test := range tests {
			fmt.Printf("test %q\n", test.text)
			scanner, err := lexer.Scanner([]byte(test.text))
			if err != nil {
				t.Fatal(err)
			}

			found := 0
			tok, err, eos := scanner.Next()
			for ; !eos; tok, err, eos = scanner.Next() {
				if err != nil {
					t.Error(err)
					fmt.Printf("err: %v\n", err)
					scanner.TC++
				} else {
					token := tok.(*Token)
					fmt.Printf("%-15v | %-30q | %d-%d | %v:%v-%v:%v\n",
						tokens[token.Type],
						strings.TrimSpace(string(token.Lexeme)),
						token.TC,
						token.TC+len(token.Lexeme),
						token.StartLine,
						token.StartColumn,
						token.EndLine,
						token.EndColumn)
					found++
				}
			}
			if found != test.tokens {
				t.Errorf("expected %v tokens got %v: %q", test.tokens, found, test.text)
			}
		}
	}
	{
		lexer := newLexer()
		if err := lexer.CompileNFA(); err != nil {
			t.Fatal(err)
		}
		runTest(lexer)
	}
	{
		lexer := newLexer()
		if err := lexer.CompileDFA(); err != nil {
			t.Fatal(err)
		}
		runTest(lexer)
	}
}

func TestNoEmptyStrings(t *testing.T) {
	skip := func(*Scanner, *machines.Match) (interface{}, error) {
		return nil, nil
	}
	lexer := NewLexer()
	lexer.Add([]byte("(ab|a)*"), skip)
	{
		if err := lexer.CompileNFA(); err == nil {
			t.Fatal("expected error")
		} else {
			t.Logf("got expected error: %v", err)
		}
	}
	{
		if err := lexer.CompileDFA(); err == nil {
			t.Fatal("expected error")
		} else {
			t.Logf("got expected error: %v", err)
		}
	}
}
