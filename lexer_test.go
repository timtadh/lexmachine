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

func TestScannerErrorToken(t *testing.T) {
	const (
		TOKEN_IDENT = iota
		TOKEN_EQUALS
		TOKEN_NUMBER
		TOKEN_PLUS
		TOKEN_ERROR
	)

	newTestLexer := func() *Lexer {
		l := NewLexer()
		l.Add([]byte("[a-zA-Z_][a-zA-Z0-9_]*"), func(s *Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(TOKEN_IDENT, string(m.Bytes), m), nil
		})
		l.Add([]byte("="), func(s *Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(TOKEN_EQUALS, nil, m), nil
		})
		l.Add([]byte("[0-9]+"), func(s *Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(TOKEN_NUMBER, string(m.Bytes), m), nil
		})
		l.Add([]byte("\\+"), func(s *Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(TOKEN_PLUS, nil, m), nil
		})
		l.Add([]byte("( |\t|\n|\r)+"), func(s *Scanner, m *machines.Match) (interface{}, error) {
			return nil, nil // skip whitespace
		})
		return l
	}

	runTests := func(compile func(*Lexer) error, engineName string) {
		// Test 1: Invalid syntax in middle of line
		{
			lexer := newTestLexer()
			if err := compile(lexer); err != nil {
				t.Fatalf("[%v] compile error: %v", engineName, err)
			}
			text := []byte("val = $ + 10")
			// '$' is at offset 6, row 1, col 7
			scanner, err := lexer.Scanner(text)
			if err != nil {
				t.Fatalf("[%v] scanner error: %v", engineName, err)
			}

			// Token 1: val
			tok, err, eos := scanner.Next()
			if err != nil || eos {
				t.Fatalf("[%v] expected valid token, got err=%v eos=%v", engineName, err, eos)
			}
			if tok.(*Token).Type != TOKEN_IDENT {
				t.Fatalf("[%v] expected IDENT, got %v", engineName, tok)
			}

			// Token 2: =
			tok, err, eos = scanner.Next()
			if err != nil || eos {
				t.Fatalf("[%v] expected valid token, got err=%v eos=%v", engineName, err, eos)
			}
			if tok.(*Token).Type != TOKEN_EQUALS {
				t.Fatalf("[%v] expected EQUALS, got %v", engineName, tok)
			}

			// Token 3: error token on '$'
			tok, err, eos = scanner.Next()
			if err == nil {
				t.Fatalf("[%v] expected error on '$', got nil", engineName)
			}
			ui, ok := err.(*machines.UnconsumedInput)
			if !ok {
				t.Fatalf("[%v] expected *machines.UnconsumedInput, got %T", engineName, err)
			}
			if ui.StartTC != 6 || ui.StartLine != 1 || ui.StartColumn != 7 {
				t.Errorf("[%v] UnconsumedInput start mismatch: got (%d, %d:%d), expected (6, 1:7)",
					engineName, ui.StartTC, ui.StartLine, ui.StartColumn)
			}
			if ui.FailTC != 7 || ui.FailLine != 1 || ui.FailColumn != 7 {
				t.Errorf("[%v] UnconsumedInput fail mismatch: got (%d, %d:%d), expected (7, 1:7)",
					engineName, ui.FailTC, ui.FailLine, ui.FailColumn)
			}

			errTok, ok := tok.(*Token)
			if !ok {
				t.Fatalf("[%v] expected error token *Token from Next(), got %T", engineName, tok)
			}
			if errTok.TC != 6 || errTok.StartLine != 1 || errTok.StartColumn != 7 {
				t.Errorf("[%v] error token start mismatch: got (%d, %d:%d), expected (6, 1:7)",
					engineName, errTok.TC, errTok.StartLine, errTok.StartColumn)
			}
			if errTok.EndLine != 1 || errTok.EndColumn != 7 {
				t.Errorf("[%v] error token end mismatch: got (%d:%d), expected (1:7)",
					engineName, errTok.EndLine, errTok.EndColumn)
			}
			if string(errTok.Lexeme) != "$" {
				t.Errorf("[%v] expected lexeme '$', got %q", engineName, string(errTok.Lexeme))
			}

			// Verify Scanner position getters
			if scanner.StartLine() != 1 || scanner.StartColumn() != 7 {
				t.Errorf("[%v] scanner start position mismatch: got (%d:%d), expected (1:7)",
					engineName, scanner.StartLine(), scanner.StartColumn())
			}
			if scanner.EndLine() != 1 || scanner.EndColumn() != 7 {
				t.Errorf("[%v] scanner end position mismatch: got (%d:%d), expected (1:7)",
					engineName, scanner.EndLine(), scanner.EndColumn())
			}

			// Verify ErrorToken helper
			customErrTok := scanner.ErrorToken(TOKEN_ERROR, err)
			if customErrTok.Type != TOKEN_ERROR || customErrTok.StartLine != 1 || customErrTok.StartColumn != 7 {
				t.Errorf("[%v] custom ErrorToken mismatch: got %v", engineName, customErrTok)
			}

			// Skip error and continue scanning
			scanner.TC = ui.FailTC

			// Token 4: +
			tok, err, eos = scanner.Next()
			if err != nil || eos {
				t.Fatalf("[%v] expected valid token after skipping error, got err=%v eos=%v", engineName, err, eos)
			}
			if tok.(*Token).Type != TOKEN_PLUS {
				t.Fatalf("[%v] expected PLUS, got %v", engineName, tok)
			}

			// Token 5: 10
			tok, err, eos = scanner.Next()
			if err != nil || eos {
				t.Fatalf("[%v] expected valid token, got err=%v eos=%v", engineName, err, eos)
			}
			if tok.(*Token).Type != TOKEN_NUMBER {
				t.Fatalf("[%v] expected NUMBER, got %v", engineName, tok)
			}
		}

		// Test 2: Multiline invalid syntax with indentation
		{
			lexer := newTestLexer()
			if err := compile(lexer); err != nil {
				t.Fatalf("[%v] compile error: %v", engineName, err)
			}
			text := []byte("line1 = 1\nline2 = 2\n  @bad\nline3 = 3")
			// "line1 = 1\n" -> 10 bytes
			// "line2 = 2\n" -> 10 bytes (offset 10 to 19)
			// "  @bad\n" -> ' ' at 20, ' ' at 21, '@' at offset 22, row 3, col 3
			scanner, err := lexer.Scanner(text)
			if err != nil {
				t.Fatalf("[%v] scanner error: %v", engineName, err)
			}

			// Consume line1 tokens: line1, =, 1
			for i := 0; i < 3; i++ {
				_, err, eos := scanner.Next()
				if err != nil || eos {
					t.Fatalf("[%v] unexpected error on line 1: %v", engineName, err)
				}
			}

			// Consume line2 tokens: line2, =, 2
			for i := 0; i < 3; i++ {
				_, err, eos := scanner.Next()
				if err != nil || eos {
					t.Fatalf("[%v] unexpected error on line 2: %v", engineName, err)
				}
			}

			// Line 3: '@' invalid syntax at row 3, col 3, byte offset 22
			tok, err, eos := scanner.Next()
			if err == nil {
				t.Fatalf("[%v] expected error on '@', got nil", engineName)
			}
			if eos {
				t.Fatalf("[%v] unexpected eos on line 3", engineName)
			}
			ui, ok := err.(*machines.UnconsumedInput)
			if !ok {
				t.Fatalf("[%v] expected *machines.UnconsumedInput, got %T", engineName, err)
			}
			if ui.StartTC != 22 || ui.StartLine != 3 || ui.StartColumn != 3 {
				t.Errorf("[%v] multiline UnconsumedInput start mismatch: got (%d, %d:%d), expected (22, 3:3)",
					engineName, ui.StartTC, ui.StartLine, ui.StartColumn)
			}
			if ui.FailTC != 23 || ui.FailLine != 3 || ui.FailColumn != 3 {
				t.Errorf("[%v] multiline UnconsumedInput fail mismatch: got (%d, %d:%d), expected (23, 3:3)",
					engineName, ui.FailTC, ui.FailLine, ui.FailColumn)
			}

			errTok := tok.(*Token)
			if errTok.TC != 22 || errTok.StartLine != 3 || errTok.StartColumn != 3 {
				t.Errorf("[%v] multiline error token mismatch: got (%d, %d:%d), expected (22, 3:3)",
					engineName, errTok.TC, errTok.StartLine, errTok.StartColumn)
			}
			if string(errTok.Lexeme) != "@" {
				t.Errorf("[%v] expected lexeme '@', got %q", engineName, string(errTok.Lexeme))
			}
		}

		// Test 3: Invalid character at end of input
		{
			lexer := newTestLexer()
			if err := compile(lexer); err != nil {
				t.Fatalf("[%v] compile error: %v", engineName, err)
			}
			text := []byte("foo#")
			// '#' at offset 3, row 1, col 4
			scanner, err := lexer.Scanner(text)
			if err != nil {
				t.Fatalf("[%v] scanner error: %v", engineName, err)
			}

			tok, err, eos := scanner.Next()
			if err != nil || eos || tok.(*Token).Type != TOKEN_IDENT {
				t.Fatalf("[%v] expected foo, got tok=%v err=%v", engineName, tok, err)
			}

			tok, err, eos = scanner.Next()
			if err == nil {
				t.Fatalf("[%v] expected error on '#', got nil", engineName)
			}
			ui := err.(*machines.UnconsumedInput)
			if ui.StartTC != 3 || ui.StartLine != 1 || ui.StartColumn != 4 {
				t.Errorf("[%v] end-of-input error start mismatch: got (%d, %d:%d), expected (3, 1:4)",
					engineName, ui.StartTC, ui.StartLine, ui.StartColumn)
			}
			if ui.FailTC != 4 || ui.FailLine != 1 || ui.FailColumn != 4 {
				t.Errorf("[%v] end-of-input error fail mismatch: got (%d, %d:%d), expected (4, 1:4)",
					engineName, ui.FailTC, ui.FailLine, ui.FailColumn)
			}
			errTok := tok.(*Token)
			if errTok.TC != 3 || errTok.StartLine != 1 || errTok.StartColumn != 4 {
				t.Errorf("[%v] end-of-input error token mismatch: got %v", engineName, errTok)
			}
		}
	}

	t.Log("Testing DFA engine")
	runTests(func(l *Lexer) error { return l.CompileDFA() }, "DFA")

	t.Log("Testing NFA engine")
	runTests(func(l *Lexer) error { return l.CompileNFA() }, "NFA")
}
