package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/timtadh/getopt"
	"github.com/timtadh/lexmachine"
)

func main() {
	short := "h"
	long := []string{
		"help",
	}

	_, optargs, err := getopt.GetOpt(os.Args[1:], short, long)
	if err != nil {
		log.Print(err)
		log.Println("try --help")
		os.Exit(1)
	}

	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			fmt.Println("parse a sensors.conf")
			os.Exit(0)
		}
	}
	lexer := newLexer()
	stmts, err := parse(lexer, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, stmt := range stmts {
		fmt.Println(stmt)
	}
}

func parse(lexer *lexmachine.Lexer, fin io.Reader) (stmts []*Node, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case error:
				err = e.(error)
				stmts = nil
			default:
				panic(e)
			}
		}
	}()
	text, err := ioutil.ReadAll(fin)
	if err != nil {
		return nil, err
	}
	scanner, err := newGoLex(lexer, text)
	if err != nil {
		return nil, err
	}
	yyParse(scanner)
	return scanner.stmts, nil
}
