package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/timtadh/getopt"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

var tokens = []string{
	"AT", "PLUS", "STAR", "DASH", "SLASH", "BACKSLASH", "CARROT", "BACKTICK", "COMMA", "LPAREN", "RPAREN",
	"BUS", "COMPUTE", "CHIP", "IGNORE", "LABEL", "SET", "NUMBER", "NAME",
	"COMMENT", "SPACE",
}
var tokmap map[string]int
var lexer *lexmachine.Lexer

func init() {
	tokmap = make(map[string]int)
	for id, name := range tokens {
		tokmap[name] = id
	}
}

func newLexer(dfa bool) *lexmachine.Lexer {
	getToken := func(tokenType int) lexmachine.Action {
		return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return s.Token(tokenType, string(m.Bytes), m), nil
		}
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
	lexer.Add([]byte(`#[^\n]*`), getToken(tokmap["COMMENT"]))
	lexer.Add([]byte(`\s+`), getToken(tokmap["SPACE"]))
	var err error
	if dfa {
		err = lexer.CompileDFA()
	} else {
		err = lexer.CompileNFA()
	}
	if err != nil {
		panic(err)
	}
	return lexer
}

func scan(text []byte) error {
	scanner, err := lexer.Scanner(text)
	if err != nil {
		return err
	}
	for tk, err, eof := scanner.Next(); !eof; tk, err, eof = scanner.Next() {
		if ui, is := err.(*machines.UnconsumedInput); ui != nil && is {
			scanner.TC = ui.FailTC
			log.Printf("skipping %v", ui)
		} else if err != nil {
			return err
		} else {
			if false {
				fmt.Println(tk)
			}
		}
	}
	return nil
}

func main() {
	short := "hdn"
	long := []string{
		"help",
		"dfa",
		"nfa",
	}

	_, optargs, err := getopt.GetOpt(os.Args[1:], short, long)
	if err != nil {
		log.Print(err)
		log.Println("help")
		os.Exit(1)
	}

	dfa := false
	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			fmt.Println("Tokenizes the standard input 1000 times")
			fmt.Println("Must supply either --nfa or --dfa. try cat /etc/sensors*.conf | sensors --nfa")
			os.Exit(0)
		case "-d", "--dfa":
			dfa = true
		case "-n", "--nfa":
			dfa = false
		}
	}
	lexer = newLexer(dfa)

	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		err = scan(text)
		if err != nil {
			log.Fatal(err)
		}
	}
}
