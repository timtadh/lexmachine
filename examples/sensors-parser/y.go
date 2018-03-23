//line sensors.y:2
package main

import __yyfmt__ "fmt"

//line sensors.y:3
import (
	"github.com/timtadh/lexmachine"
)

//line sensors.y:11
type yySymType struct {
	yys   int
	token *lexmachine.Token
	ast   *Node
}

const AT = 57346
const PLUS = 57347
const STAR = 57348
const DASH = 57349
const SLASH = 57350
const BACKSLASH = 57351
const CARROT = 57352
const BACKTICK = 57353
const COMMA = 57354
const LPAREN = 57355
const RPAREN = 57356
const BUS = 57357
const COMPUTE = 57358
const CHIP = 57359
const IGNORE = 57360
const LABEL = 57361
const SET = 57362
const NUMBER = 57363
const NAME = 57364
const NEWLINE = 57365

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"AT",
	"PLUS",
	"STAR",
	"DASH",
	"SLASH",
	"BACKSLASH",
	"CARROT",
	"BACKTICK",
	"COMMA",
	"LPAREN",
	"RPAREN",
	"BUS",
	"COMPUTE",
	"CHIP",
	"IGNORE",
	"LABEL",
	"SET",
	"NUMBER",
	"NAME",
	"NEWLINE",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line sensors.y:101

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 64

var yyAct = [...]int{

	31, 29, 35, 18, 30, 11, 14, 12, 15, 13,
	16, 38, 20, 4, 32, 38, 41, 34, 33, 28,
	39, 21, 26, 25, 39, 24, 23, 40, 37, 36,
	22, 19, 37, 36, 27, 47, 48, 49, 45, 10,
	46, 50, 9, 43, 51, 44, 54, 55, 52, 53,
	8, 43, 56, 44, 43, 2, 44, 17, 42, 7,
	6, 5, 3, 1,
}
var yyPact = [...]int{

	-10, -10, -1000, -20, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 9, -1, 8, 4, 3, 1, -1000, -1000, 0,
	-1000, -1, -3, 7, -1000, 7, -6, -1000, -1000, 46,
	32, -1000, 11, 11, 11, -1000, -1000, -1000, -1000, 7,
	49, -1000, 7, 7, 7, 7, 7, -1000, -1000, -1000,
	38, 49, 32, 32, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 63, 55, 62, 61, 60, 59, 50, 42, 39,
	12, 1, 4, 0, 2,
}
var yyR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	3, 10, 10, 4, 5, 6, 7, 8, 9, 11,
	11, 11, 12, 12, 12, 13, 13, 13, 13, 14,
	14, 14, 14,
}
var yyR2 = [...]int{

	0, 2, 1, 2, 1, 1, 1, 1, 1, 1,
	1, 2, 1, 4, 2, 3, 5, 2, 3, 3,
	3, 1, 3, 3, 1, 2, 2, 2, 1, 1,
	1, 1, 3,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, 23, -4, -5, -6, -7, -8,
	-9, 15, 17, 19, 16, 18, 20, -2, 23, 22,
	-10, 22, 22, 22, 22, 22, 22, -10, 22, -11,
	-12, -13, 7, 11, 10, -14, 22, 21, 4, 13,
	-11, 22, 12, 5, 7, 6, 8, -14, -14, -14,
	-11, -11, -12, -12, -13, -13, 14,
}
var yyDef = [...]int{

	0, -2, 2, 0, 4, 5, 6, 7, 8, 9,
	10, 0, 0, 0, 0, 0, 0, 1, 3, 0,
	14, 12, 0, 0, 17, 0, 0, 11, 15, 0,
	21, 24, 0, 0, 0, 28, 29, 30, 31, 0,
	18, 13, 0, 0, 0, 0, 0, 25, 26, 27,
	0, 16, 19, 20, 22, 23, 32,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 3:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sensors.y:43
		{
			yylex.(*golex).stmts = append(yylex.(*golex).stmts, yyDollar[1].ast)
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:47
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:48
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:49
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:50
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:51
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:52
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sensors.y:55
		{
			yyVAL.ast = yyDollar[2].ast.PrependKid(NewNode("name", yyDollar[1].token))
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:56
		{
			yyVAL.ast = NewNode("names", nil).AddKid(NewNode("name", yyDollar[1].token))
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sensors.y:59
		{
			yyVAL.ast = NewNode("bus", yyDollar[1].token).AddKid(NewNode("name", yyDollar[2].token)).AddKid(NewNode("name", yyDollar[3].token)).AddKid(NewNode("name", yyDollar[4].token))
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sensors.y:62
		{
			yyVAL.ast = NewNode("chip", yyDollar[1].token).AddKid(yyDollar[2].ast)
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sensors.y:65
		{
			yyVAL.ast = NewNode("label", yyDollar[1].token).AddKid(NewNode("name", yyDollar[2].token)).AddKid(NewNode("name", yyDollar[3].token))
		}
	case 16:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sensors.y:69
		{
			yyVAL.ast = NewNode("compute", yyDollar[1].token).AddKid(NewNode("name", yyDollar[2].token)).AddKid(yyDollar[3].ast).AddKid(yyDollar[5].ast)
		}
	case 17:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sensors.y:72
		{
			yyVAL.ast = NewNode("ignore", yyDollar[1].token).AddKid(NewNode("name", yyDollar[2].token))
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sensors.y:75
		{
			yyVAL.ast = NewNode("set", yyDollar[1].token).AddKid(NewNode("name", yyDollar[2].token)).AddKid(yyDollar[3].ast)
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sensors.y:78
		{
			yyVAL.ast = NewNode("+", yyDollar[2].token).AddKid(yyDollar[1].ast).AddKid(yyDollar[3].ast)
		}
	case 20:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sensors.y:79
		{
			yyVAL.ast = NewNode("-", yyDollar[2].token).AddKid(yyDollar[1].ast).AddKid(yyDollar[3].ast)
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:80
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sensors.y:83
		{
			yyVAL.ast = NewNode("*", yyDollar[2].token).AddKid(yyDollar[1].ast).AddKid(yyDollar[3].ast)
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sensors.y:84
		{
			yyVAL.ast = NewNode("/", yyDollar[2].token).AddKid(yyDollar[1].ast).AddKid(yyDollar[3].ast)
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:85
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sensors.y:88
		{
			yyVAL.ast = NewNode("negate", yyDollar[1].token).AddKid(yyDollar[2].ast)
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sensors.y:89
		{
			yyVAL.ast = NewNode("`", yyDollar[1].token).AddKid(yyDollar[2].ast)
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sensors.y:90
		{
			yyVAL.ast = NewNode("^", yyDollar[1].token).AddKid(yyDollar[2].ast)
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:91
		{
			yyVAL.ast = yyDollar[1].ast
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:94
		{
			yyVAL.ast = NewNode("name", yyDollar[1].token)
		}
	case 30:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:95
		{
			yyVAL.ast = NewNode("number", yyDollar[1].token)
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sensors.y:96
		{
			yyVAL.ast = NewNode("@", yyDollar[1].token)
		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sensors.y:97
		{
			yyVAL.ast = yyDollar[2].ast
		}
	}
	goto yystack /* stack new state and value */
}
