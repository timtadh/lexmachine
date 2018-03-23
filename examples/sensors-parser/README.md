# Parse a sensors.conf file with goyacc and lexmachine

This is an example for how to integrate lexmachine with the [standard
yacc](http://godoc.org/golang.org/x/tools/cmd/goyacc) implementation for Go.
Yacc is its own weird and interesting language for specifying bottom up shift
reduce parsers. You can "easily" use lexmachine with yacc but it does require
some understanding of

1.  How yacc works (eg. the things it generates)
2.  How to use those generated definitions in your code

## Running the example

```sh
$ go generate -x -v github.com/timtadh/lexmachine/examples/sensors-parser
$ go install github.com/timtadh/lexmachine/examples/sensors-parser
$ cat examples/sensors-parser/sensors.conf | sensors-parser
```

## Partial Explanation

Yacc controls the definitions for Tokens with its `%token` directives (see the
`sensors.y` file. You will use those definitions in your lexer. An example of
how to do this is in `sensors_golex.go`.

Second, Yacc expects the lexer to conform to the following interface:

```go
type yyLexer interface {
   // Lex gets the next token and puts it in lval
   Lex(lval *yySymType) (tokenType int)
   // Error is called on parse error (it should probably panic?)
   Error(message string)
 }
```

The `yySymType` is generate from Yacc via the `%union` directive. The tokenType
is the token identifier. However, the tokenType needs to be in the correct range
for Yacc which starts at `yyPrivate`. The way to get the types identified
correctly is to set it as `return token.Type + yyPrivate - 1`. See
`sensors_golex.go` for a full example.

Yacc in its own special way has each production "return" a yySymType which
serves as *both* an AST node *and* a token. Thus, my definition for yySymType
is:

```yacc
%union{
    token *lexmachine.Token
    ast   *Node
}
```

This lets you construct an AST while parsing:

```yacc
Unary : DASH Factor             { $$.ast = NewNode("negate", $1.token).AddKid($2.ast) }
      | BACKTICK Factor         { $$.ast = NewNode("`", $1.token).AddKid($2.ast) }
      | CARROT Factor           { $$.ast = NewNode("^", $1.token).AddKid($2.ast) }
      | Factor                  { $$.ast = $1.ast }
      ;

Factor : NAME                   { $$.ast = NewNode("name", $1.token) }
       | NUMBER                 { $$.ast = NewNode("number", $1.token) }
       | AT                     { $$.ast = NewNode("@", $1.token) }
       | LPAREN Expr RPAREN     { $$.ast = $2.ast }
       ;
```

Finally, yacc does not provide any means of returning anything from the parser.
To deal with this the lexer you provide needs to have a field which provides the
result back to the caller:

```go
type golex struct {
    *lexmachine.Scanner
    stmts []*Node
}
```

In the example `stmts` provides the parsed statements from the file back to the
caller of the parser:

```go
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
```

Since `yyLexer` is an interface in goyacc their is some casting involved to
populate `stmts` (note `yylex` is a magic variable in yacc that refers to the
lexer object you provided):

```yacc
Line : Stmt NEWLINE             { yylex.(*golex).stmts = append(yylex.(*golex).stmts, $1.ast) }
     | NEWLINE
     ;

```
