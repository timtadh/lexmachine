package frontend

import (
	"testing"

	"github.com/timtadh/data-structures/test"
)

func followEquals(follow []map[int]bool, expected [][]int) bool {
	if len(follow) != len(expected) {
		return false
	}
	for row, cols := range expected {
		if len(cols) != len(follow[row]) {
			return false
		}
		for _, c := range cols {
			if !follow[row][c] {
				return false
			}
		}
	}
	return true
}

func TestFollowExample(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(a|b)*xyz"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
		NewCharacter('x'),
		NewCharacter('y'),
		NewCharacter('z'),
	}
	expectedFollows := [][]int{
		{0, 1, 2},
		{0, 1, 2},
		{3},
		{4},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowStar(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a*b"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
	}
	expectedFollows := [][]int{
		{0, 1},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowPlus(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a+b"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
	}
	expectedFollows := [][]int{
		{0, 1},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowMaybe(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("ab?c"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
		NewCharacter('c'),
	}
	expectedFollows := [][]int{
		{1, 2},
		{2},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowMaybes(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("ab?c?(d|e|f)?g?h"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'), // 0
		NewCharacter('b'), // 1
		NewCharacter('c'), // 2
		NewCharacter('d'), // 3
		NewCharacter('e'), // 4
		NewCharacter('f'), // 5
		NewCharacter('g'), // 6
		NewCharacter('h'), // 7
	}
	expectedFollows := [][]int{
		{1, 2, 3, 4, 5, 6, 7},
		{2, 3, 4, 5, 6, 7},
		{3, 4, 5, 6, 7},
		{6, 7},
		{6, 7},
		{6, 7},
		{7},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowMaybeStar(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(ab?c?)*d"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
		NewCharacter('c'),
		NewCharacter('d'),
	}
	expectedFollows := [][]int{
		{0, 1, 2, 3},
		{0, 2, 3},
		{0, 3},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowMaybeNested(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a(b?c)?d"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
		NewCharacter('c'),
		NewCharacter('d'),
	}
	expectedFollows := [][]int{
		{1, 2, 3},
		{2},
		{3},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowMaybeNested2(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a(bc?)?d"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
		NewCharacter('c'),
		NewCharacter('d'),
	}
	expectedFollows := [][]int{
		{1, 3},
		{2, 3},
		{3},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowMaybeNested3(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("q((bc?|x|y)?)*z"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('q'), // 0
		NewCharacter('b'), // 1
		NewCharacter('c'), // 2
		NewCharacter('x'), // 3
		NewCharacter('y'), // 4
		NewCharacter('z'), // 5
	}
	expectedFollows := [][]int{
		{1, 3, 4, 5},
		{2, 5, 1, 3, 4},
		{5, 1, 3, 4},
		{5, 1, 3, 4},
		{5, 1, 3, 4},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowMaybeNested4(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a(((bc?)?d?)?(ef?))?g"))
	t.AssertNil(err)
	program, err := Generate(ast)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t_match(program, "abeg", t)
	expectedPos := []AST{
		NewCharacter('a'), // 0
		NewCharacter('b'), // 1
		NewCharacter('c'), // 2
		NewCharacter('d'), // 3
		NewCharacter('e'), // 4
		NewCharacter('f'), // 5
		NewCharacter('g'), // 6
	}
	expectedFollows := [][]int{
		{1, 3, 4, 6},
		{2, 3, 4},
		{3, 4},
		{4},
		{5, 6},
		{6},
		{},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

func TestFollowNested(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("(a+b)*(c|d)*"))
	t.AssertNil(err)
	expectedPos := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
		NewCharacter('c'),
		NewCharacter('d'),
	}
	expectedFollows := [][]int{
		{0, 1},
		{0, 2, 3},
		{2, 3},
		{2, 3},
	}
	positions, follow := Follow(ast)
	t.Assert(listEquals(expectedPos, positions), "follow \n\tproduced: %v\n\texpected: %v", positions, expectedPos)
	t.Assert(followEquals(follow, expectedFollows), "follow \n\tproduced: %v\n\texpected: %v", follow, expectedFollows)
}

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

func listEquals(a, b []AST) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Equals(b[i]) {
			return false
		}
	}
	return true
}

func TestFirst_char(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a"))
	t.AssertNil(err)
	first := []AST{
		NewCharacter('a'),
	}
	t.Assert(listEquals(first, ast.First()), "first \n\tproduced: %v, \n\texpected: %v", ast.First(), first)
}

func TestLast_char(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a"))
	t.AssertNil(err)
	last := []AST{
		NewCharacter('a'),
	}
	t.Assert(listEquals(last, ast.Last()), "last \n\tproduced: %v\n\texpected: %v", ast.Last(), last)
}

func TestFirst_range(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[a-z]"))
	t.AssertNil(err)
	first := []AST{
		NewRange('a', 'z'),
	}
	t.Assert(listEquals(first, ast.First()), "first \n\tproduced: %v, \n\texpected: %v", ast.First(), first)
}

func TestLast_range(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("[a-z]"))
	t.AssertNil(err)
	last := []AST{
		NewRange('a', 'z'),
	}
	t.Assert(listEquals(last, ast.Last()), "last \n\tproduced: %v\n\texpected: %v", ast.Last(), last)
}

func TestFirst_ops(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a*+?"))
	t.AssertNil(err)
	first := []AST{
		NewCharacter('a'),
	}
	t.Assert(listEquals(first, ast.First()), "first \n\tproduced: %v, \n\texpected: %v", ast.First(), first)
}

func TestLast_ops(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a*+?"))
	t.AssertNil(err)
	last := []AST{
		NewCharacter('a'),
	}
	t.Assert(listEquals(last, ast.Last()), "last \n\tproduced: %v\n\texpected: %v", ast.Last(), last)
}

func TestFirst_alt(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a|b"))
	t.AssertNil(err)
	first := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
	}
	t.Assert(listEquals(first, ast.First()), "first \n\tproduced: %v, \n\texpected: %v", ast.First(), first)
}

func TestLast_alt(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a|b"))
	t.AssertNil(err)
	last := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
	}
	t.Assert(listEquals(last, ast.Last()), "last \n\tproduced: %v\n\texpected: %v", ast.Last(), last)
}

func TestFirst_concat(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("a?b?c?de"))
	t.AssertNil(err)
	first := []AST{
		NewCharacter('a'),
		NewCharacter('b'),
		NewCharacter('c'),
		NewCharacter('d'),
	}
	t.Assert(listEquals(first, ast.First()), "first \n\tproduced: %v, \n\texpected: %v", ast.First(), first)
}

func TestLast_concat(x *testing.T) {
	t := (*test.T)(x)
	ast, err := Parse([]byte("abc?d?e?"))
	t.AssertNil(err)
	last := []AST{
		NewCharacter('e'),
		NewCharacter('d'),
		NewCharacter('c'),
		NewCharacter('b'),
	}
	t.Assert(listEquals(last, ast.Last()), "last \n\tproduced: %v\n\texpected: %v", ast.Last(), last)
}
