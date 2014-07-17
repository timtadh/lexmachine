package frontend

import (
	"fmt"
)

var (
	ErrorNoOp = fmt.Errorf("No Operator")
)

func match_any(text []byte, i int) (int, AST, error) {
	if i >= len(text) {
		return i, nil, fmt.Errorf("out of text, %d", i)
	}
	return i+1, NewCharacter(text[i]), nil
}

func match(text []byte, i int, c byte) (int, error) {
	if i >= len(text) {
		return i, fmt.Errorf("out of text, %d", i)
	} else if text[i] == c {
		i++
		return i, nil
	}
	return i,
	fmt.Errorf(
		"Expected text at pos, %d, to equal %s, got %s",
		i,
		string([]byte{c}),
		string(text[i:i+1]),
	)
}

func Parse(text []byte) (AST, error) {
	return regex(text)
}

func regex(text []byte) (AST, error) {
	i, ast, err := alternation(text, 0)
	if err != nil {
		return nil, err
	} else if i != len(text) {
		return nil, fmt.Errorf("unconsumed input")
	}
	return NewMatch(ast), nil
}

func alternation(text []byte, i int) (int, AST, error) {
	return _alt(text, i)
}

func _alt(text []byte, i int) (int, AST, error) {
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

func alternation_(text []byte, i int) (int, AST, error) {
	if i >= len(text) {
		return i, nil, nil
	}
	i, err := match(text, i, '|')
	if err != nil {
		return i, nil, nil
	}
	return _alt(text, i)
}

func choice(text []byte, i int) (int, AST, error) {
	return _choice(text, i)
}

func _choice(text []byte, i int) (int, AST, error) {
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

func choice_(text []byte, i int) (int, AST, error) {
	if i >= len(text) {
		return i, nil, nil
	}
	i, C, _ := _choice(text, i)
	return i, C, nil
}

func atomicOp(text []byte, i int) (int, AST, error) {
	i, A, err := atomic(text, i)
	if err != nil {
		return i, nil, err
	}
	i, O, err := op(text, i)
	if err != nil && err == ErrorNoOp {
		return i, A, nil
	} else if err != nil {
		return i, A, err
	}
	return i, NewApplyOp(O, A), err
}

func op(text []byte, i int) (int, AST, error) {
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
	return i, nil, ErrorNoOp
}

func atomic(text []byte, i int) (int, AST, error) {
	i, ast, err := char(text, i)
	if err == nil {
		return i, ast, nil
	}
	i, ast, err = group(text, i)
	if err == nil {
		return i, ast, nil
	}
	return i, nil, fmt.Errorf("Expected group or concat at %d", i)
}

func group(text []byte, i int) (int, AST, error) {
	i, err := match(text, i, '(')
	if err != nil {
		return i, nil, err
	}
	i, A, err := alternation(text, i)
	if err != nil {
		return i, nil, err
	}
	i, err = match(text, i, ')')
	if err != nil {
		return i, nil, err
	}
	return i, A, nil
}

func concat(text []byte, i int) (int, AST, error) {
	return _concat(text, i)
}

func _concat(text []byte, i int) (int, AST, error) {
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

func concat_(text []byte, i int) (int, AST, error) {
	if i >= len(text) {
		return i, nil, nil
	}
	i, C, _ := _concat(text, i)
	return i, C, nil
}

func char(text []byte, i int) (int, AST, error) {
	i, C, err := CHAR(text, i)
	if err == nil {
		return i, C, nil
	}
	i, R, err := charRange(text, i)
	if err == nil {
		return i, R, nil
	}
	return i, nil, fmt.Errorf(
		"Expected a CHAR or charRange at %d", i)
}

func CHAR(text []byte, i int) (int, AST, error) {
	i, err := match(text, i, '\\')
	if err == nil {
		if i < len(text) && text[i] == 'n' {
			return i+1, NewCharacter('\n'), nil
		}
		return match_any(text, i)
	}
	if i >= len(text) {
		return i, nil, fmt.Errorf(
			"ran out of text in CHAR, %d", i)
	}
	switch text[i] {
	case '|','+','*','?','(',')','[',']', '^':
		return i, nil, fmt.Errorf(
			"unexpected operator, %s", string([]byte{text[i]}))
	case '.':
		return i+1, NewAny(), nil
	default:
		return match_any(text, i)
	}
}

func charRange(text []byte, i int) (int, AST, error) {
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

func charNotRange(text []byte, i int) (int, AST, error) {
	if i >= len(text) {
		return i, nil, fmt.Errorf("out of text, %d", i)
	}
	ch := i
	i++
	i, err := match(text, i, ']')
	if err != nil {
		return i, nil, err
	}
	ast := NewAlternation(
		NewRange(NewCharacter(0), NewCharacter(text[ch]-1)),
		NewRange(NewCharacter(text[ch]+1), NewCharacter(255)),
	)
	return i, ast, nil
}

