# `lexmachine` - Lexical Analysis Framework for Golang

By Tim Henderson

Copyright 2015, All Rights Reserved. Made available to for public use
under the terms of the GNU General Public License version 3 or any later
version (at your option).

## What?

`lexmachine` is a full lexical analysis framework for the Go programming
language. It supports a restricted but usable set of regular expressions
appropriate for writing lexers for complex programming languages. The
framework also supports sub lexers and non-regular lexing through an
"escape hatch" which allows the users to consume any number of further
bytes after a match. So if you want to support nested C-style comments
or other paired structures you can do so at the lexical analysis stage.

## Documentation

- [Tutorial](http://hackthology.com/writing-a-lexer-in-go-with-lexmachine.html)
- [![GoDoc](https://godoc.org/github.com/timtadh/lexmachine?status.svg)](https://godoc.org/github.com/timtadh/lexmachine)

## History

This library was written when I was teaching EECS 337 *Compiler Design
and Implementation* at Case Western Reserve University in Fall of 2014.
It wrote two compilers one was "hidden" from the students  as the
language implemented was their project language. The other was
[tcel](https://github.com/timtadh/tcel) which was written initially as
an example of how to do type checking. That compiler was later expanded
to explain AST interpretation, intermediate code generation, and x86
code generation.

## What is in Box

`lexmachine` includes the following components

1. A parser for restricted set of regular expressions.
2. A abstract syntax tree (AST) for regular expressions.
3. A backpatching code generator which compiles the AST to (NFA) machine
   code.
4. An NFA simulation based lexical analysis engine. Lexical analysis
   engines work in a slightly different way from a normal regular
   expression engine as they tokenize a stream rather than matching one
   string.
5. Match objects which include start and end column and line numbers of
   the lexemes as well as their associate token name.
6. A declarative "DSL" for specifying the lexers.
7. An "escape hatch" which allows one to match non-regular tokens by
   consuming any number of further bytes after the match.
8. An in progress (but working) DFA backend which current runs DFA
   simulation but is intended to support code generation as well.
9. A command `lexc` which compiles a sequence of patterns into an NFA.
   Mostly written to support a homework assignment for the class.

I am just now starting to document this project. My intention is to get
it cleaned up with examples and get the DFA backend finished up and
merged. I also want to rewrite the frontend to support the full set of
regular expressions. I wasn't able to do that initially because of time
constraints.

