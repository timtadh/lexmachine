package frontend

import (
	"fmt"
	"sort"
	"strings"
	"runtime"
)

import (
	"github.com/timtadh/data-structures/errors"
)

var (
	DEBUG = false
)

type ParseError struct {
	Reason string
	Production string
	TC int
	text []byte
	chain []*ParseError
}

func Errorf(text []byte, tc int, format string, args ...interface{}) *ParseError {
	pc, _, _, ok := runtime.Caller(1)
	var fn string = "unknown"
	if ok {
		fn = runtime.FuncForPC(pc).Name()
		split := strings.Split(fn, ".")
		fn = split[len(split)-1]
	}
	msg := fmt.Sprintf(format, args...)
	return &ParseError{
		Reason: msg,
		Production: fn,
		TC: tc,
		text: text,
	}
}

func linecol(text []byte, tc int) (line int, col int) {
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

func (p *ParseError) Error() string {
	errs := make([]string, 0, len(p.chain) + 1)
	for i := len(p.chain)-1; i >= 0; i-- {
		errs = append(errs, p.chain[i].Error())
	}
	errs = append(errs, p.error())
	return strings.Join(errs, "\n")
}

func (p *ParseError) error() string {
	line, col := linecol(p.text, p.TC)
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

func match_any(text []byte, i int) (int, AST, *ParseError) {
	if i >= len(text) {
		return i, nil, Errorf(text, i, "out of text, %d", i)
	}
	return i+1, NewCharacter(text[i]), nil
}

func match(text []byte, i int, c byte) (int, *ParseError) {
	if i >= len(text) {
		return i, Errorf(text, i, "out of text, %d", i)
	} else if text[i] == c {
		i++
		return i, nil
	}
	return i,
	Errorf(text, i, 
		"Expected text at pos, %d, to equal %s, got %s",
		i,
		string([]byte{c}),
		string(text[i:i+1]),
	)
}

func Parse(text []byte) (AST, error) {
	a, err := regex(text)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func regex(text []byte) (AST, *ParseError) {
	i, ast, err := alternation(text, 0)
	if err != nil {
		return nil, err
	} else if i != len(text) {
		return nil, Errorf(text, i, "unconsumed input")
	}
	return NewMatch(ast), nil
}

func alternation(text []byte, i int) (int, AST, *ParseError) {
	return _alt(text, i)
}

func _alt(text []byte, i int) (int, AST, *ParseError) {
	i, C, err := choice(text, i)
	if err != nil {
		return i, nil, err
	}
	i, A, err := alternation_(text, i)
	if err != nil {
		return i, nil, err
	}
	return i, NewAlternation(C, A), nil
}

func alternation_(text []byte, i int) (int, AST, *ParseError) {
	if i >= len(text) {
		return i, nil, nil
	}
	i, err := match(text, i, '|')
	if err != nil {
		return i, nil, nil
	}
	return _alt(text, i)
}

func choice(text []byte, i int) (int, AST, *ParseError) {
	return _choice(text, i)
}

func _choice(text []byte, i int) (int, AST, *ParseError) {
	i, A, err := atomicOp(text, i)
	if err != nil {
		return i, nil, err
	}
	i, C, err := choice_(text, i)
	if err != nil {
		return i, nil, err
	}
	return i, NewConcat(A, C), nil
}

func choice_(text []byte, i int) (int, AST, *ParseError) {
	if i >= len(text) {
		return i, nil, nil
	}
	i, C, _ := _choice(text, i)
	return i, C, nil
}

func atomicOp(text []byte, i int) (int, AST, *ParseError) {
	i, A, err := atomic(text, i)
	if DEBUG {
		errors.Logf("DEBUG", "atomic %v", err)
	}
	if err != nil {
		return i, nil, err
	}
	i, O, err := op(text, i)
	if err != nil && err.Reason == "No Operator" {
		return i, A, nil
	} else if err != nil {
		return i, A, err
	}
	return i, NewApplyOp(O, A), err
}

func op(text []byte, i int) (int, AST, *ParseError) {
	i, err := match(text, i, '+')
	if err == nil {
		return i, NewOp("+"), nil
	}
	i, err = match(text, i, '*')
	if err == nil {
		return i, NewOp("*"), nil
	}
	i, err = match(text, i, '?')
	if err == nil {
		return i, NewOp("?"), nil
	}
	return i, nil, Errorf(text, i, "No Operator")
}

func atomic(text []byte, i int) (int, AST, *ParseError) {
	i, ast, errChar := char(text, i)
	if errChar == nil {
		return i, ast, nil
	}
	if DEBUG {
		errors.Logf("DEBUG", "char %v", errChar)
	}
	i, ast, errGroup := group(text, i)
	if errGroup == nil {
		return i, ast, nil
	}
	if DEBUG {
		errors.Logf("DEBUG", "group %v", errGroup)
	}
	return i, nil, Errorf(text, i, "Expected group or char").Chain(errChar).Chain(errGroup)
}

func group(text []byte, j int) (int, AST, *ParseError) {
	i, err := match(text, j, '(')
	if err != nil {
		return i, nil, err
	}
	i, A, err := alternation(text, i)
	if err != nil {
		return j, nil, err
	}
	i, err = match(text, i, ')')
	if err != nil {
		return j, nil, err
	}
	return i, A, nil
}

func concat(text []byte, i int) (int, AST, *ParseError) {
	return _concat(text, i)
}

func _concat(text []byte, i int) (int, AST, *ParseError) {
	i, Ch, err := char(text, i)
	if err != nil {
		return i, nil, err
	}
	i, Co, err := concat_(text, i)
	if err != nil {
		return i, nil, err
	}
	return i, NewConcat(Ch, Co), nil
}

func concat_(text []byte, i int) (int, AST, *ParseError) {
	if i >= len(text) {
		return i, nil, nil
	}
	return _concat(text, i)
}

func char(text []byte, i int) (int, AST, *ParseError) {
	i, C, errCHAR := CHAR(text, i)
	if errCHAR == nil {
		return i, C, nil
	}
	i, R, errRange := charRange(text, i)
	if errRange == nil {
		return i, R, nil
	}
	return i, nil, Errorf(text, i,
		"Expected a CHAR or charRange at %d, %v", i, string(text)).Chain(errCHAR).Chain(errRange)
}

func CHAR(text []byte, i int) (int, AST, *ParseError) {
	if i >= len(text) {
		return i, nil, Errorf(text, i, "out of input %v, %v", i, string(text))
	}
	if text[i] == '\\' {
		i, b, err := getByte(text, i)
		if err != nil {
			return i, nil, err
		}
		return i+1, NewCharacter(b), nil
	}
	switch text[i] {
	case '|','+','*','?','(',')','[',']', '^':
		return i, nil, Errorf(text, i, 
			"unexpected operator, %s", string([]byte{text[i]}))
	case '.':
		return i+1, NewAny(), nil
	default:
		return i+1, NewCharacter(text[i]), nil
	}
}

func getByte(text []byte, i int) (int, byte, *ParseError) {
	i, err := match(text, i, '\\')
	if err == nil {
		if i < len(text) && text[i] == 'n' {
			return i, '\n', nil
		} else if i < len(text) && text[i] == 'r' {
			return i, '\r', nil
		} else if i < len(text) && text[i] == 't' {
			return i, '\t', nil
		}
		return i, text[i], nil
	}
	if i >= len(text) {
		return i, 0, Errorf(text, i, "ran out of text at %d", i)
	}
	return i, text[i], nil
}

func charRange(text []byte, i int) (int, AST, *ParseError) {
	i, err := match(text, i, '[')
	if err != nil {
		return i, nil, err
	}
	i, err = match(text, i, '^')
	if err == nil {
		return charNotRange(text, i)
	}
	i, S, err := match_any(text, i)
	if err != nil {
		return i, nil, err
	}
	i, err = match(text, i, '-')
	if err != nil {
		return i, nil, err
	}
	i, T, err := match_any(text, i)
	if err != nil {
		return i, nil, err
	}
	i, err = match(text, i, ']')
	if err != nil {
		return i, nil, err
	}
	return i, NewRange(S, T), err
}

func charNotRange(text []byte, i int) (int, AST, *ParseError) {
	if i >= len(text) {
		return i, nil, Errorf(text, i, "out of text, %d", i)
	}
	chs := make([]byte, 0, 10)
	for ; i < len(text) && text[i] != ']'; i++ {
		var b byte
		var err *ParseError
		i, b, err = getByte(text, i)
		if err != nil {
			return i, nil, err
		}
		chs = append(chs, b)
	}
	i, err := match(text, i, ']')
	if err != nil {
		return i, nil, err
	}
	if len(chs) == 0 {
		return i, nil, Errorf(text, i, "empty negated range at %v", i)
	}
	sortBytes(chs)
	ranges := make([]*Range, 0, len(chs)+1)
	var prev byte = 0
	for _, ch := range chs {
		if prev == ch {
			goto loop_inc
		}
		ranges = append(ranges, &Range{From: prev, To: ch-1})
		loop_inc:
			prev = ch+1
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


