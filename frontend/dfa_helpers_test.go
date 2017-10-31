package frontend

import "testing"
import "github.com/timtadh/data-structures/test"

func TestMatchesEmptyString_char(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a"))
	t.AssertNil(err)
	t.Assert(!ast.MatchesEmptyString(), "character should not match the empty string, %v", ast)
}

func TestMatchesEmptyString_range(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[a-z]"))
	t.AssertNil(err)
	t.Assert(!ast.MatchesEmptyString(), "range should not match the empty string, %v", ast)
}

func TestMatchesEmptyString_maybe(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a?"))
	t.AssertNil(err)
	t.Assert(ast.MatchesEmptyString(), "maybe should match the empty string, %v", ast)
}

func TestMatchesEmptyString_star(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a*"))
	t.AssertNil(err)
	t.Assert(ast.MatchesEmptyString(), "star should match the empty string, %v", ast)
}

func TestMatchesEmptyString_plus(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a+"))
	t.AssertNil(err)
	t.Assert(!ast.MatchesEmptyString(), "a+ should not match the empty string, %v", ast)
	ast, err = Parse([]byte("a?+"))
	t.AssertNil(err)
	t.Assert(ast.MatchesEmptyString(), "a?+ should match the empty string, %v", ast)
}

func TestMatchesEmptyString_alt(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a|b"))
	t.AssertNil(err)
	t.Assert(!ast.MatchesEmptyString(), "a|b should not match the empty string, %v", ast)
	ast, err = Parse([]byte("a?|b"))
	t.AssertNil(err)
	t.Assert(ast.MatchesEmptyString(), "a?|b should match the empty string, %v", ast)
	ast, err = Parse([]byte("a|b?"))
	t.AssertNil(err)
	t.Assert(ast.MatchesEmptyString(), "a|b? should match the empty string, %v", ast)
	ast, err = Parse([]byte("a?|b?"))
	t.AssertNil(err)
	t.Assert(ast.MatchesEmptyString(), "a?|b? should match the empty string, %v", ast)
}

func TestMatchesEmptyString_altMatch(x *testing.T) {
	t := (*test.T)(x)
	a, err := Parse([]byte("a"))
	t.AssertNil(err)
	maybe_a, err := Parse([]byte("a?"))
	t.AssertNil(err)
	b, err := Parse([]byte("b"))
	t.AssertNil(err)
	maybe_b, err := Parse([]byte("b?"))
	t.AssertNil(err)
	ast := NewAltMatch(a, b)
	t.Assert(!ast.MatchesEmptyString(), "a|b should not match the empty string, %v", ast)
	ast = NewAltMatch(maybe_a, b)
	t.Assert(ast.MatchesEmptyString(), "a?|b should match the empty string, %v", ast)
	ast = NewAltMatch(a, maybe_b)
	t.Assert(ast.MatchesEmptyString(), "a|b? should match the empty string, %v", ast)
	ast = NewAltMatch(maybe_a, maybe_b)
	t.Assert(ast.MatchesEmptyString(), "a?|b? should match the empty string, %v", ast)
}

func TestMatchesEmptyString_concat(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("ab"))
	t.AssertNil(err)
	t.Assert(!ast.MatchesEmptyString(), "ab should not match the empty string, %v", ast)
	ast, err = Parse([]byte("a?b"))
	t.AssertNil(err)
	t.Assert(!ast.MatchesEmptyString(), "a?b should not match the empty string, %v", ast)
	ast, err = Parse([]byte("ab?"))
	t.AssertNil(err)
	t.Assert(!ast.MatchesEmptyString(), "ab? should not match the empty string, %v", ast)
	ast, err = Parse([]byte("a?b?"))
	t.AssertNil(err)
	t.Assert(ast.MatchesEmptyString(), "a?b? should match the empty string, %v", ast)
}
