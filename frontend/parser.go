package frontend

import (
	"fmt"
	"log"
	"runtime"
	"sort"
	"strings"
)

var (
	DEBUG = false
)

type ParseError struct {
	Reason     string
	Production string
	TC         int
	text       []byte
	chain      []*ParseError
}

func Errorf(text []byte, tc int, format string, args ...interface{}) *ParseError {
	pc, _, _, ok := runtime.Caller(1)
	return errorf(pc, ok, text, tc, format, args...)
}

func matchErrorf(text []byte, tc int, format string, args ...interface{}) *ParseError {
	pc, _, _, ok := runtime.Caller(2)
	return errorf(pc, ok, text, tc, format, args...)
}

func errorf(pc uintptr, ok bool, text []byte, tc int, format string, args ...interface{}) *ParseError {
	var fn string = "unknown"
	if ok {
		fn = runtime.FuncForPC(pc).Name()
		split := strings.Split(fn, ".")
		fn = split[len(split)-1]
	}
	msg := fmt.Sprintf(format, args...)
	return &ParseError{
		Reason:     msg,
		Production: fn,
		TC:         tc,
		text:       text,
	}
}

func (p *ParseError) Error() string {
	errs := make([]string, 0, len(p.chain)+1)
	for i := len(p.chain) - 1; i >= 0; i-- {
		errs = append(errs, p.chain[i].Error())
	}
	errs = append(errs, p.error())
	return strings.Join(errs, "\n")
}

func (p *ParseError) error() string {
	line, col := LineCol(p.text, p.TC)
	return fmt.Sprintf("Regex parse error in production '%v' : at index %v line %v column %v '%s' : %v",
		p.Production, p.TC, line, col, p.text[p.TC:], p.Reason)
}

func (p *ParseError) String() string {
	return p.Error()
}

func (p *ParseError) Chain(e *ParseError) *ParseError {
	p.chain = append(p.chain, e)
	return p
}

func LineCol(text []byte, tc int) (line int, col int) {
	for i := 0; i <= tc && i < len(text); i++ {
		if text[i] == '\n' {
			col = 0
			line += 1
		} else {
			col += 1
		}
	}
	if tc == 0 && tc < len(text) {
		if text[tc] == '\n' {
			line += 1
			col -= 1
		}
	}
	return line, col
}

func Parse(text []byte) (AST, error) {
	a, err := (&parser{
		text:      text,
		lastError: Errorf(text, 0, "unconsumed input"),
	}).regex()
	if err != nil {
		return nil, err
	}
	return a, nil
}

type parser struct {
	text      []byte
	lastError *ParseError
}

func (p *parser) regex() (AST, *ParseError) {
	i, ast, err := p.alternation(0)
	if err != nil {
		return nil, err
	} else if i != len(p.text) {
		return nil, p.lastError
	}
	return NewMatch(ast), nil
}

func (p *parser) alternation(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter alternation %v '%v'", i, string(p.text[i:]))
	}
	i, A, err := p.atomicOps(i)
	if err != nil {
		return i, nil, err
	}
	i, B, err := p.alternation_(i)
	if err != nil {
		return i, nil, err
	}
	return i, NewAlternation(A, B), nil
}

func (p *parser) alternation_(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter alternation_ %v '%v'", i, string(p.text[i:]))
	}
	if i >= len(p.text) {
		return i, nil, nil
	}
	i, err := p.match(i, '|')
	if err != nil {
		return i, nil, nil
	}
	i, A, err := p.atomicOps(i)
	if err != nil {
		return i, nil, err
	}
	i, B, err := p.alternation_(i)
	if err != nil {
		return i, nil, err
	}
	return i, NewAlternation(A, B), nil
}

func (p *parser) atomicOps(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter atomicOps %v '%v'", i, string(p.text[i:]))
	}
	if i >= len(p.text) {
		return i, nil, nil
	}
	i, A, err := p.atomicOp(i)
	if err != nil {
		p.lastError.Chain(err)
		return i, nil, nil
	}
	i, B, err := p.atomicOps(i)
	if err != nil {
		return i, nil, err
	}
	return i, NewConcat(A, B), nil
}

func (p *parser) atomicOp(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter atomicOp %v '%v'", i, string(p.text[i:]))
	}
	i, A, err := p.atomic(i)
	if DEBUG {
		log.Printf("atomic %v", err)
	}
	if err != nil {
		return i, nil, err
	}
	i, O, err := p.op(i)
	if err != nil && err.Reason == "No Operator" {
		return i, A, nil
	} else if err != nil {
		return i, A, err
	}
	return i, NewApplyOp(O, A), err
}

func (p *parser) op(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter op %v '%v'", i, string(p.text[i:]))
	}
	i, err := p.match(i, '+')
	if err == nil {
		return i, NewOp("+"), nil
	}
	i, err = p.match(i, '*')
	if err == nil {
		return i, NewOp("*"), nil
	}
	i, err = p.match(i, '?')
	if err == nil {
		return i, NewOp("?"), nil
	}
	return i, nil, Errorf(p.text, i, "No Operator")
}

func (p *parser) atomic(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter atomic %v '%v'", i, string(p.text[i:]))
	}
	i, ast, errChar := p.char(i)
	if errChar == nil {
		return i, ast, nil
	}
	if DEBUG {
		log.Printf("char %v", errChar)
	}
	i, ast, errGroup := p.group(i)
	if errGroup == nil {
		return i, ast, nil
	}
	if DEBUG {
		log.Printf("group %v", errGroup)
	}
	return i, nil, Errorf(p.text, i, "Expected group or char").Chain(errChar).Chain(errGroup)
}

func (p *parser) group(j int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter group %v '%v'", j, string(p.text[j:]))
	}
	i, err := p.match(j, '(')
	if err != nil {
		return i, nil, err
	}
	i, A, err := p.alternation(i)
	if err != nil {
		return j, nil, err
	}
	i, err = p.match(i, ')')
	if err != nil {
		return j, nil, err
	}
	return i, A, nil
}

func (p *parser) concat(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter concat %v '%v'", i, string(p.text[i:]))
	}
	if i >= len(p.text) {
		return i, nil, nil
	}
	i, Ch, err := p.char(i)
	if err != nil {
		return i, nil, nil
	}
	i, Co, err := p.concat(i)
	if err != nil {
		return i, Ch, nil
	}
	return i, NewConcat(Ch, Co), nil
}

func (p *parser) char(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter char %v '%v'", i, string(p.text[i:]))
	}
	i, C, errCHAR := p.CHAR(i)
	if errCHAR == nil {
		return i, C, nil
	}
	i, R, errRange := p.charRange(i)
	if errRange == nil {
		return i, R, nil
	}
	return i, nil, Errorf(p.text, i,
		"Expected a CHAR or charRange at %d, %v", i, string(p.text)).Chain(errCHAR).Chain(errRange)
}

func (p *parser) CHAR(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter CHAR %v '%v'", i, string(p.text[i:]))
	}
	if i >= len(p.text) {
		return i, nil, Errorf(p.text, i, "out of input %v, %v", i, string(p.text))
	}
	if p.text[i] == '\\' {
		i, b, err := p.getByte(i)
		if err != nil {
			return i, nil, err
		}
		return i + 1, NewCharacter(b), nil
	}
	switch p.text[i] {
	case '|', '+', '*', '?', '(', ')', '[', ']', '^':
		return i, nil, Errorf(p.text, i,
			"unexpected operator, %s", string([]byte{p.text[i]}))
	case '.':
		return i + 1, NewAny(), nil
	default:
		return i + 1, NewCharacter(p.text[i]), nil
	}
}

func (p *parser) getByte(i int) (int, byte, *ParseError) {
	i, err := p.match(i, '\\')
	if err == nil {
		if i < len(p.text) && p.text[i] == 'n' {
			return i, '\n', nil
		} else if i < len(p.text) && p.text[i] == 'r' {
			return i, '\r', nil
		} else if i < len(p.text) && p.text[i] == 't' {
			return i, '\t', nil
		}
		return i, p.text[i], nil
	}
	if i >= len(p.text) {
		return i, 0, Errorf(p.text, i, "ran out of p.text at %d", i)
	}
	return i, p.text[i], nil
}

func (p *parser) charRange(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter charRange %v '%v'", i, string(p.text[i:]))
	}
	i, err := p.match(i, '[')
	if err != nil {
		return i, nil, err
	}
	i, err = p.match(i, '^')
	if err == nil {
		return p.charNotRange(i)
	}
	i, S, err := p.match_any(i)
	if err != nil {
		return i, nil, err
	}
	i, err = p.match(i, '-')
	if err != nil {
		return i, nil, err
	}
	i, T, err := p.match_any(i)
	if err != nil {
		return i, nil, err
	}
	i, err = p.match(i, ']')
	if err != nil {
		return i, nil, err
	}
	return i, NewRange(S, T), err
}

func (p *parser) charNotRange(i int) (int, AST, *ParseError) {
	if DEBUG {
		log.Printf("enter charNotRange %v '%v'", i, string(p.text[i:]))
	}
	if i >= len(p.text) {
		return i, nil, Errorf(p.text, i, "out of p.text, %d", i)
	}
	chs := make([]byte, 0, 10)
	for ; i < len(p.text) && p.text[i] != ']'; i++ {
		var b byte
		var err *ParseError
		i, b, err = p.getByte(i)
		if err != nil {
			return i, nil, err
		}
		chs = append(chs, b)
	}
	i, err := p.match(i, ']')
	if err != nil {
		return i, nil, err
	}
	if len(chs) == 0 {
		return i, nil, Errorf(p.text, i, "empty negated range at %v", i)
	}
	sortBytes(chs)
	ranges := make([]*Range, 0, len(chs)+1)
	var prev byte = 0
	for _, ch := range chs {
		if prev == ch {
			goto loop_inc
		}
		ranges = append(ranges, &Range{From: prev, To: ch - 1})
	loop_inc:
		prev = ch + 1
	}
	ast := NewAlternation(
		ranges[len(ranges)-1],
		&Range{From: prev, To: 255},
	)
	for j := len(ranges) - 2; j >= 0; j-- {
		ast = NewAlternation(
			ranges[j],
			ast,
		)
	}
	return i, ast, nil
}

func (p *parser) match_any(i int) (int, AST, *ParseError) {
	if i >= len(p.text) {
		return i, nil, Errorf(p.text, i, "out of p.text, %d", i)
	}
	return i + 1, NewCharacter(p.text[i]), nil
}

func (p *parser) match(i int, c byte) (int, *ParseError) {
	if i >= len(p.text) {
		return i, matchErrorf(p.text, i, "out of p.text, %d", i)
	} else if p.text[i] == c {
		i++
		return i, nil
	}
	return i,
		matchErrorf(p.text, i,
			"expected '%v' at %v got '%v' of '%v'",
			string([]byte{c}),
			i,
			string(p.text[i:i+1]),
			string(p.text[i:]),
		)
}

func sortBytes(bytes sortableBytes) {
	sort.Sort(bytes)
}

type sortableBytes []byte

func (s sortableBytes) Len() int {
	return len(s)
}

func (s sortableBytes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortableBytes) Less(i, j int) bool {
	return s[i] < s[j]
}
