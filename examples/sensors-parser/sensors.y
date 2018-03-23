%{

package main

import (
    "github.com/timtadh/lexmachine"
)

%}

%union{
    token *lexmachine.Token
    ast   *Node
}

%token	AT
%token	PLUS
%token	STAR
%token	DASH
%token	SLASH
%token	BACKSLASH
%token	CARROT
%token	BACKTICK
%token	COMMA
%token	LPAREN
%token	RPAREN
%token	BUS
%token	COMPUTE
%token	CHIP
%token	IGNORE
%token	LABEL
%token	SET
%token	NUMBER
%token	NAME
%token	NEWLINE

%% /* The grammar follows.  */

Lines : Lines Line
      | Line
      ;

Line : Stmt NEWLINE             { yylex.(*golex).stmts = append(yylex.(*golex).stmts, $1.ast) }
     | NEWLINE
     ;

Stmt : Bus                      { $$.ast = $1.ast }
     | Chip                     { $$.ast = $1.ast }
     | Label                    { $$.ast = $1.ast }
     | Compute                  { $$.ast = $1.ast }
     | Ignore                   { $$.ast = $1.ast }
     | Set                      { $$.ast = $1.ast }
     ;

Names : NAME Names              { $$.ast = $2.ast.PrependKid(NewNode("name", $1.token)) }
      | NAME                    { $$.ast = NewNode("names", nil).AddKid(NewNode("name", $1.token)) }
      ;

Bus : BUS NAME NAME NAME        { $$.ast = NewNode("bus", $1.token).AddKid(NewNode("name", $2.token)).AddKid(NewNode("name", $3.token)).AddKid(NewNode("name", $4.token)) }
    ;

Chip : CHIP Names               { $$.ast = NewNode("chip", $1.token).AddKid($2.ast) }
     ;

Label : LABEL NAME NAME         { $$.ast = NewNode("label", $1.token).AddKid(NewNode("name", $2.token)).AddKid(NewNode("name", $3.token)) }
      ;

Compute : COMPUTE NAME Expr COMMA Expr
                                { $$.ast = NewNode("compute", $1.token).AddKid(NewNode("name", $2.token)).AddKid($3.ast).AddKid($5.ast)  }
        ;

Ignore : IGNORE NAME            { $$.ast = NewNode("ignore", $1.token).AddKid(NewNode("name", $2.token)) }
       ;

Set : SET NAME Expr             { $$.ast = NewNode("set", $1.token).AddKid(NewNode("name", $2.token)).AddKid($3.ast) }
    ;

Expr : Expr PLUS Term           { $$.ast = NewNode("+", $2.token).AddKid($1.ast).AddKid($3.ast) }
     | Expr DASH Term           { $$.ast = NewNode("-", $2.token).AddKid($1.ast).AddKid($3.ast) }
     | Term                     { $$.ast = $1.ast }
     ;

Term : Term STAR Unary          { $$.ast = NewNode("*", $2.token).AddKid($1.ast).AddKid($3.ast) }
     | Term SLASH Unary         { $$.ast = NewNode("/", $2.token).AddKid($1.ast).AddKid($3.ast) }
     | Unary                    { $$.ast = $1.ast }
     ;

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

;
%%
