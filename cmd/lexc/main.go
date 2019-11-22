package main

import (
	"fmt"
	logpkg "log"
	"os"
)

import (
	"github.com/timtadh/getopt"
)

import (
	"github.com/timtadh/lexmachine/frontend"
)

var log *logpkg.Logger

func init() {
	log = logpkg.New(os.Stderr, "", 0)
}

var usageMessage = "lexc -p <pattern> [-p <pattern>]*"
var extendedMessage = `
lexc compiles regular expressions to a program suitable for lexing

Options
    -h, --help                          print this message
    -p, --pattern=<pattern>             a regex pattern

Specs
    <pattern>
        a regex pattern
`

func usage(code int) {
	fmt.Fprintln(os.Stderr, usageMessage)
	if code == 0 {
		fmt.Fprintln(os.Stderr, extendedMessage)
		code = 1
	} else {
		fmt.Fprintln(os.Stderr, "Try -h or --help for help")
	}
	os.Exit(code)
}

func main() {

	short := "hp:"
	long := []string{
		"help",
		"pattern=",
	}

	_, optargs, err := getopt.GetOpt(os.Args[1:], short, long)
	if err != nil {
		log.Print(err)
		usage(1)
	}

	patterns := make([]string, 0, 10)
	for _, oa := range optargs {
		switch oa.Opt() {
		case "-h", "--help":
			usage(0)
		case "-p", "--pattern":
			patterns = append(patterns, oa.Arg())
		}
	}

	if len(patterns) <= 0 {
		log.Print("Must supply some regulars expressions!")
		usage(1)
	}

	asts := make([]frontend.AST, 0, len(patterns))
	for _, p := range patterns {
		ast, err := frontend.Parse([]byte(p))
		if err != nil {
			log.Fatal(err)
		}
		asts = append(asts, ast)
	}

	lexast := asts[len(asts)-1]
	for i := len(asts) - 2; i >= 0; i-- {
		lexast = frontend.NewAltMatch(asts[i], lexast)
	}

	program, err := frontend.Generate(lexast)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(program.Serialize())
}
