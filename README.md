# `lexmachine` - Lexical Analysis Framework for Golang

By Tim Henderson

Copyright 2014-2017, All Rights Reserved. Made available for public use under
the terms of a BSD 3-Clause license.

[![GoDoc](https://godoc.org/github.com/timtadh/lexmachine?status.svg)](https://godoc.org/github.com/timtadh/lexmachine)
[![ReportCard](https://goreportcard.com/badge/github.com/timtadh/lexmachine)](https://goreportcard.com/report/github.com/timtadh/lexmachine)

## What?

`lexmachine` is a full lexical analysis framework for the Go programming
language. It supports a restricted but usable set of regular expressions
appropriate for writing lexers for complex programming languages. The framework
also supports sub lexers and non-regular lexing through an "escape hatch" which
allows the users to consume any number of further bytes after a match. So if you
want to support nested C-style comments or other paired structures you can do so
at the lexical analysis stage.

## Goal

`lexmachine` intends to be the best, fastest, and easiest to use lexical
analysis system for Go.

## Documentation

-   [Tutorial](http://hackthology.com/writing-a-lexer-in-go-with-lexmachine.html)
-   [![GoDoc](https://godoc.org/github.com/timtadh/lexmachine?status.svg)](https://godoc.org/github.com/timtadh/lexmachine)

## What is in Box

`lexmachine` includes the following components

1.  A parser for restricted set of regular expressions.
2.  A abstract syntax tree (AST) for regular expressions.
3.  A backpatching code generator which compiles the AST to (NFA) machine code.
4.  Both DFA (Deterministic Finite Automata) and a NFA (Non-deterministic Finite
    Automata) simulation based lexical analysis engines. Lexical analysis
    engines work in a slightly different way from a normal regular expression
    engine as they tokenize a stream rather than matching one string.
5.  Match objects which include start and end column and line numbers of the
    lexemes as well as their associate token name.
6.  A declarative "DSL" for specifying the lexers.
7.  An "escape hatch" which allows one to match non-regular tokens by consuming
    any number of further bytes after the match.
8.  A command `lexc` which compiles a sequence of patterns into an NFA. Mostly
    written to support a homework assignment for the class.

## Quick Start

Import lexmachine

```go
import (
    "github.com/timtadh/lexmachine"
    "github.com/timtadh/lexmachine/machines"
)
```

Define your tokens and types

```go
var Literals []string // The tokens representing literal strings
var Keywords []string // The keyword tokens
var Tokens []string // All of the tokens (including literals and keywords)
var TokenIds map[string]int // A map from the token names to their int ids
var Lexer *lex.Lexer // The lexer object. Use this to construct a Scanner

func init() {
    Literals = []string{
        "[",
        "]",
        "{",
        "}",
        "=",
        ",",
        ";",
        ":",
        "->",
        "--",
    }
    Keywords = []string{
        "NODE",
        "EDGE",
        "GRAPH",
        "DIGRAPH",
        "SUBGRAPH",
        "STRICT",
    }
    Tokens = []string{
        "COMMENT",
        "ID",
    }
    Tokens = append(Tokens, Keywords...)
    Tokens = append(Tokens, Literals...)
    TokenIds = make(map[string]int)
    for i, tok := range Tokens {
        TokenIds[tok] = i
    }
}
```

Create a lexer object and add few regular expression patterns to match

```go

// token produces a new token (of the supplied name) constructed from a match
func token(name string) lex.Action {
    return func(s *lex.Scanner, m *machines.Match) (interface{}, error) {
        return s.Token(TokenIds[name], string(m.Bytes), m), nil
    }
}

// CreateLexer defines a lexer for the graphviz dot language.
func CreateLexer() (*lexmachine.Lexer, error) {
    lexer := lexmachine.NewLexer()

    for _, lit := range Literals {
        r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
        lexer.Add([]byte(r), token(lit))
    }
    for _, name := range Keywords {
        lexer.Add([]byte(strings.ToLower(name)), token(name))
    }

    lexer.Add([]byte(`//[^\n]*\n?`), token("COMMENT"))
    lexer.Add([]byte(`/\*([^*]|\r|\n|(\*+([^*/]|\r|\n)))*\*+/`), token("COMMENT"))
    lexer.Add([]byte(`([a-z]|[A-Z])([a-z]|[A-Z]|[0-9]|_)*`), token("ID"))
    lexer.Add([]byte(`"([^\\"]|(\\.))*"`), token("ID"))
    lexer.Add([]byte("( |\t|\n|\r)+"), skip)
    lexer.Add([]byte(`\<`),
        func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
            str := make([]byte, 0, 10)
            str = append(str, match.Bytes...)
            brackets := 1
            match.EndLine = match.StartLine
            match.EndColumn = match.StartColumn
            for tc := scan.TC; tc < len(scan.Text); tc++ {
                str = append(str, scan.Text[tc])
                match.EndColumn += 1
                if scan.Text[tc] == '\n' {
                    match.EndLine += 1
                }
                if scan.Text[tc] == '<' {
                    brackets += 1
                } else if scan.Text[tc] == '>' {
                    brackets -= 1
                }
                if brackets == 0 {
                    match.TC = scan.TC
                    scan.TC = tc + 1
                    match.Bytes = str
                    return token("ID")(scan, match)
                }
            }
            return nil,
                fmt.Errorf("unclosed HTML literal starting at %d, (%d, %d)",
                    match.TC, match.StartLine, match.StartColumn)
        },
    )

    // Compile a DFA you can also use CompileNFA() to compile an NFA instead.
    err := lexer.CompileDFA()
    if err != nil {
        return nil, err
    }
    return lexer, nil
}
```

Scan a string to tokenize it.

```go
func ExampleLex() error {
    lexer, err := CreateLexer()
    if err != nil {
        return err
    }
    scanner, err := lexer.Scanner([]byte(`digraph {
      rankdir=LR;
      a [label="a" shape=box];
      c [<label>=<<u>C</u>>];
      b [label="bb"];
      a -> c;
      c -> b;
      d -> c;
      b -> a;
      b -> e;
      e -> f;
    }`))
    if err != nil {
        return err
    }
    fmt.Println("Type    | Lexeme     | Position")
    fmt.Println("--------+------------+------------")
    for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
        if ui, is := err.(*machines.UnconsumedInput); is {
            // if you wanted to skip bad tokens you could
            // scanner.TC = ui.FailTC
            // here we are just going to fail on a bad token
            return ui
        } else if err != nil {
            return err
        }
        token := tok.(*lexmachine.Token)
        fmt.Printf("%-7v | %-10v | %v:%v-%v:%v\n",
            dot.Tokens[token.Type],
            string(token.Lexeme),
            token.StartLine,
            token.StartColumn,
            token.EndLine,
            token.EndColumn)
    }
    return nil
}
```

## Regular Expressions

Lexmachine (like most lexical analysis frameworks) uses [Regular
Expressions](https://en.wikipedia.org/wiki/Regular_expression) to specify the
*patterns* to match when spitting the string up into categorized *tokens.* 
For a more advanced introduction to regular expressions engines see Russ Cox's
[articles](https://swtch.com/~rsc/regexp/). To learn more about how regular
expressions are used to *tokenize* string take a look at Alex Aiken's [video
lectures](https://youtu.be/SRhkfvqeA1M) on the subject. Finally, Aho *et al.*
give a through treatment of the subject in the [Dragon
Book](http://www.worldcat.org/oclc/951336275) Chapter 3.

A regular expression is a *pattern* which *matches* a set of strings. It is made
up of *characters* such as `a` or `b`, characters with special meanings (such as
`.` which matches any character), and operators. The regular expression `abc`
matches exactly one string `abc`.

### Charater Expressions

In lexmachine most characters (eg. `a`, `b` or `#`) represent themselves. Some
have special meanings (as detailed below in operators). However, all characters
can be represented by prefixing the character with a `\`.

#### Any Character

`.` matches any character.

#### Special Characters

1. `\` use `\\` to match
2. newline use `\n` to match
3. cariage return use `\r` to match
4. tab use `\t` to match
5. `.` use `\.` to match
6. operators: {`|`, `+`, `*`, `?`, `(`, `)`, `[`, `]`, `^`} prefix with a `\` to
   match.

#### Character Classes

Sometimes it is advantages to match a variety of characters. For instance, if
you want to ignore captilization for the work `Capitol` you could write the
expression `[Cc]aptiol` which would match both `Capitol` or `capitol`. There are
two forms of character ranges:

1. `[abcd]` matches all the letters inside the `[]` (eg. that pattern matches
   the strings `a`, `b`, `c`, `d`).
2. `[a-d]` matches the range of characters between the character before the dash
   (`a`) and the character after the dash (`d`) (eg. that pattern matches
   the strings `a`, `b`, `c`, `d`).

These two forms may be combined:

For instance, `[a-zA-Z123]` matches the strings {`a`, `b`, ..., `z`, `A`, `B`,
... `Z`, `0`, `2`, `3`}

#### Inverted Character Classes

Sometimes it is easier to specify the characters you don't want to match than
the characters you do. For instance, you want to match any character but a lower
case one. This can be achieved using an inverted class: `[^a-z]`. An inverted
class is specified by putting a `^` just after the opening bracket.

#### Built-in Character Classes

1. `\d` = `[0-9]` (the digit class)
2. `\D` = `[^0-9]` (the not a digit class)
3. `\s` = `[ \t\n\r\f]` (the space class). where \f is a form feed (note: \f is
   not a special sequence in lexmachine, if you want to specify the form feed
   character (ascii 0x0c) use []byte{12}.
4. `\S` = `[^ \t\n\r\f]` (the not a space class)
5. `\w` = `[0-9a-zA-Z_]` (the letter class)
5. `\W` = `[^0-9a-zA-Z_]` (the not a letter class)

### Operators

1. The pipe operator `|` indicates alternative choices. For instance the
   expression `a|b` matches either the string `a` or the string `b` but not `ab`
   or `ba` or the empty string.

2. The parenthesis operator `()` groups a subexpression together. For instance
   the expression `a(b|c)d` matches `abd` or `acd` but not `abcd`.

3. The star operator `*` indicates the "starred" subexpression should match zero
   or more times. For instance, `a*` matches the empty string, `a`, `aa`, `aaa`
   and so on.

4. The plus operator `+` indicates the "plussed" subexpression should match one
   or more times. For instance, `a+` matches `a`, `aa`, `aaa` and so on.

5. The maybe operator `?` indicates the "questioned" subexpression should match
   zero or one times. For instance, `a?` matches the empty string and `a`.

### Grammar

The canonical grammar is found in the handwritten recursive descent
[parser](https://github.com/timtadh/lexmachine/blob/master/frontend/parser.go).
This section should be considered documentation not specification.

Note: e stands for the empty string

```
Regex -> Alternation

Alternation -> AtomicOps Alternation'

Alternation' -> `|` AtomicOps Alternation'
              | e

AtomicOps -> AtomicOp AtomicOps
           | e

AtomicOp -> Atomic
          | Atomic Ops

Ops -> Op Ops
     | e

Op -> `+`
    | `*`
    | `?`

Atomic -> Char
        | Group

Group -> `(` Alternation `)`

Char -> CHAR
      | CharClass

CharClass -> `[` Range `]`
           | `[` `^` Range `]`

Range -> CharClassItem Range'

Range' -> CharClassItem Range'
        | e

CharClassItem -> BYTE
              -> BYTE `-` BYTE

CHAR -> matches any character expect '|', '+', '*', '?', '(', ')', '[', ']', '^'
        unless escaped. Additionally '.' is returned a as the wildcard character
        which matches any character. Built-in character classes are also handled
        here.

BYTE -> matches any byte
```

## History

This library was started when I was teaching EECS 337 *Compiler Design and
Implementation* at Case Western Reserve University in Fall of 2014. It wrote two
compilers one was "hidden" from the students as the language implemented was
their project language. The other was [tcel](https://github.com/timtadh/tcel)
which was written initially as an example of how to do type checking. That
compiler was later expanded to explain AST interpretation, intermediate code
generation, and x86 code generation.
