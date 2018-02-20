package lexmachine

import (
	"fmt"
	"strconv"
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
	skip := func(*Scanner, *machines.Match) (interface{}, error) {
		return nil, nil
	}
	token := func(id int, name string) Action {
		return func(s *Scanner, m *machines.Match) (interface{}, error) {
			return string(m.Bytes), nil
		}
	}

	data := "true" // This input fails.
	// data := "true " // this with a trailing space does not.

	lexer := NewLexer()
	lexer.Add([]byte("true"), token(0, "TRUE"))
	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	if err := lexer.CompileDFA(); err != nil {
		t.Fatal(err)
	}

	var scanner *Scanner

	scanner, err := lexer.Scanner([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	found := 0
	tok, err, eos := scanner.Next()
	for ; !eos; tok, err, eos = scanner.Next() {
		fmt.Printf("Token: %v\n", tok)
		found++
	}
	if found != 1 {
		t.Errorf("Expected exactly 1 tokens got %v, ===\nErr: %v\nEOS: %v\nTC: %d\n", found, err, eos, scanner.TC)

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

# Use this to send dhcp log messages to a different log file (you also
# have to hack syslog.conf to complete the redirection).
#log-facility local7;

# No service will be given on this subnet, but declaring it helps the 
# DHCP server to understand the network topology.

#subnet 10.152.187.0 netmask 255.255.255.0 {
#}

# This is a very basic subnet declaration.

#subnet 10.254.239.0 netmask 255.255.255.224 {
#  range 10.254.239.10 10.254.239.20;
#  option routers rtr-239-0-1.example.org, rtr-239-0-2.example.org;
#}

# This declaration allows BOOTP clients to get dynamic addresses,
# which we don't really recommend.

#subnet 10.254.239.32 netmask 255.255.255.224 {
#  range dynamic-bootp 10.254.239.40 10.254.239.60;
#  option broadcast-address 10.254.239.31;
#  option routers rtr-239-32-1.example.org;
#}

# A slightly different configuration for an internal subnet.
#subnet 10.5.5.0 netmask 255.255.255.224 {
#  range 10.5.5.26 10.5.5.30;
#  option domain-name-servers ns1.internal.example.org;
#  option domain-name "internal.example.org";
#  option routers 10.5.5.1;
#  option broadcast-address 10.5.5.31;
#  default-lease-time 600;
#  max-lease-time 7200;
#}

# Hosts which require special configuration options can be listed in
# host statements.   If no address is specified, the address will be
# allocated dynamically (if possible), but the host-specific information
# will still come from the host declaration.

#host passacaglia {
#  hardware ethernet 0:0:c0:5d:bd:95;
#  filename "vmunix.passacaglia";
#  server-name "toccata.example.com";
#}

# Fixed IP addresses can also be specified for hosts.   These addresses
# should not also be listed as being available for dynamic assignment.
# Hosts for which fixed IP addresses have been specified can boot using
# BOOTP or DHCP.   Hosts for which no fixed address is specified can only
# be booted with DHCP, unless there is an address range on the subnet
# to which a BOOTP client is connected which has the dynamic-bootp flag
# set.
#host fantasia {
#  hardware ethernet 08:00:07:26:c0:a5;
#  fixed-address fantasia.example.com;
#}

# You can declare a class of clients and then do address allocation
# based on that.   The example below shows a case where all clients
# in a certain class get addresses on the 10.17.224/24 subnet, and all
# other clients get addresses on the 10.0.29/24 subnet.

#class "foo" {
#  match if substring (option vendor-class-identifier, 0, 4) = "SUNW";
#}

#shared-network 224-29 {
#  subnet 10.17.224.0 netmask 255.255.255.0 {
#    option routers rtr-224.example.org;
#  }
#  subnet 10.0.29.0 netmask 255.255.255.0 {
#    option routers rtr-29.example.org;
#  }
#  pool {
#    allow members of "foo";
#    range 10.17.224.10 10.17.224.250;
#  }
#  pool {
#    deny members of "foo";
#    range 10.0.29.10 10.0.29.230;
#  }
#}
`

	newLexer := func() *Lexer {
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
		lex.Add([]byte(`([a-z]|[A-Z]|[0-9]|_|\-|\.)*`), token("ID"))
		lex.Add([]byte(`"([^\\"]|(\\.))*"`), token("ID"))
		lex.Add([]byte("[\n \t]"), skip)
		for _, lit := range literals {
			lex.Add([]byte(lit), token(lit))
		}

		err := lex.Compile()
		if err != nil {
			panic(err)
		}

		return lex
	}

	scanner, err := newLexer().Scanner([]byte(text))
	if err != nil {
		return
	}
	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {
		if err != nil {
			t.Error(err)
		}
		token := tok.(*Token)
		fmt.Printf("%-7v | %-10v | %v:%v-%v:%v\n",
			token.Type,
			string(token.Lexeme),
			token.StartLine,
			token.StartColumn,
			token.EndLine,
			token.EndColumn)
	}
}
