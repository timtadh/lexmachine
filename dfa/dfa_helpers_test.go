package dfa

import (
	"testing"

	"github.com/timtadh/data-structures/test"
	"github.com/timtadh/lexmachine/frontend"
)

func TestLabeledAst(x *testing.T) {
	t := (*test.T)(x)
	verify := func(ast frontend.AST) {
		last := Label(ast)
		pos := 0
		var visit func(int, frontend.AST) int
		visit = func(i int, n frontend.AST) int {
			for _, kid := range n.Children() {
				i = visit(i, kid)
			}
			t.Assert(last.Order[i].Equals(n), "Expected %v got %v", n, last.Order[i])
			t.Assert(len(last.Kids[i]) == len(n.Children()), "Expected %v children got %v", len(n.Children()), len(last.Kids[i]))
			for j := 0; j < len(last.Kids[i]); j++ {
				t.Assert(last.Order[last.Kids[i][j]].Equals(n.Children()[j]), "Expected %v got %v", n.Children()[j], last.Order[last.Kids[i][j]])
			}
			switch n.(type) {
			case *frontend.Character, *frontend.Range:
				t.Assert(last.Order[last.Positions[pos]].Equals(n), "Expected %v got %v", n, last.Order[last.Positions[pos]])
				pos++
			}
			return i + 1
		}
		visit(0, last.Root)
	}
	for _, regex := range []string{
		"a", "b", "asdf", "s|a", "sdf*", "(sdf)+(asdf)*", "w|(s|e)*(s)+(s?fe)**", "(a|we|f*|s*?)|W(LSD)Adf[23-s]",
	} {
		ast, err := frontend.Parse([]byte(regex))
		t.AssertNil(err)
		verify(frontend.DesugarRanges(ast))
	}
}

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

func testFollow(x *testing.T, regex string, expectedPos []frontend.AST, expectedFollows [][]int) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte(regex))
	t.AssertNil(err)
	lAst := Label(ast)
	positions := lAst.Positions
	_, follow := lAst.Follow()
	t.Assert(listEquals(expectedPos, astListOrder(lAst, positions)),
		"%v follow \n\tproduced: %v\n\texpected: %v", regex, astListOrder(lAst, positions), expectedPos)
	t.Assert(followEquals(follow, expectedFollows),
		"%v follow \n\tproduced: %v\n\texpected: %v", regex, follow, expectedFollows)
}

func TestFollowSimple(x *testing.T) {
	testFollow(
		x,
		"a*",
		[]frontend.AST{
			frontend.NewCharacter('a'),
			frontend.NewEOS(),
		},
		[][]int{
			{0, 1},
			{},
		})
}

func TestFollowExample(x *testing.T) {
	testFollow(
		x,
		"(a|b)*xyz",
		[]frontend.AST{
			frontend.NewCharacter('a'),
			frontend.NewCharacter('b'),
			frontend.NewCharacter('x'),
			frontend.NewCharacter('y'),
			frontend.NewCharacter('z'),
			frontend.NewEOS(),
		},
		[][]int{
			{0, 1, 2},
			{0, 1, 2},
			{3},
			{4},
			{5},
			{},
		})
}

func TestFollowStar(x *testing.T) {
	testFollow(
		x,
		"a*b",
		[]frontend.AST{
			frontend.NewCharacter('a'),
			frontend.NewCharacter('b'),
			frontend.NewEOS(),
		},
		[][]int{
			{0, 1},
			{2},
			{},
		})
}

func TestFollowStar2(x *testing.T) {
	testFollow(
		x,
		"ab*c",
		[]frontend.AST{
			frontend.NewCharacter('a'),
			frontend.NewCharacter('b'),
			frontend.NewCharacter('c'),
			frontend.NewEOS(),
		},
		[][]int{
			{1, 2},
			{1, 2},
			{3},
			{},
		})
}

func TestFollowPlus(x *testing.T) {
	testFollow(
		x,
		"a+b",
		[]frontend.AST{
			frontend.NewCharacter('a'),
			frontend.NewCharacter('b'),
			frontend.NewEOS(),
		},
		[][]int{
			{0, 1},
			{2},
			{},
		})
}

func TestFollowMaybe(x *testing.T) {
	testFollow(
		x,
		"ab?c",
		[]frontend.AST{
			frontend.NewCharacter('a'),
			frontend.NewCharacter('b'),
			frontend.NewCharacter('c'),
			frontend.NewEOS(),
		},
		[][]int{
			{1, 2},
			{2},
			{3},
			{},
		})
}

func TestFollowMaybes(x *testing.T) {
	testFollow(
		x,
		"ab?c?(d|e|f)?g?h",
		[]frontend.AST{
			frontend.NewCharacter('a'), // 0
			frontend.NewCharacter('b'), // 1
			frontend.NewCharacter('c'), // 2
			frontend.NewCharacter('d'), // 3
			frontend.NewCharacter('e'), // 4
			frontend.NewCharacter('f'), // 5
			frontend.NewCharacter('g'), // 6
			frontend.NewCharacter('h'), // 7
			frontend.NewEOS(),
		},
		[][]int{
			{1, 2, 3, 4, 5, 6, 7},
			{2, 3, 4, 5, 6, 7},
			{3, 4, 5, 6, 7},
			{6, 7},
			{6, 7},
			{6, 7},
			{7},
			{8},
			{},
		})
}

func TestFollowMaybeStar(x *testing.T) {
	testFollow(
		x,
		"(ab?c?)*d",
		[]frontend.AST{
			frontend.NewCharacter('a'), // 0
			frontend.NewCharacter('b'), // 1
			frontend.NewCharacter('c'), // 2
			frontend.NewCharacter('d'), // 3
			frontend.NewEOS(),
		},
		[][]int{
			{0, 1, 2, 3},
			{0, 2, 3},
			{0, 3},
			{4},
			{},
		})
}

func TestFollowMaybeNested(x *testing.T) {
	testFollow(
		x,
		"a(b?c)?d",
		[]frontend.AST{
			frontend.NewCharacter('a'), // 0
			frontend.NewCharacter('b'), // 1
			frontend.NewCharacter('c'), // 2
			frontend.NewCharacter('d'), // 3
			frontend.NewEOS(),
		},
		[][]int{
			{1, 2, 3},
			{2},
			{3},
			{4},
			{},
		})
}

func TestFollowMaybeNested2(x *testing.T) {
	testFollow(
		x,
		"a(bc?)?d",
		[]frontend.AST{
			frontend.NewCharacter('a'), // 0
			frontend.NewCharacter('b'), // 1
			frontend.NewCharacter('c'), // 2
			frontend.NewCharacter('d'), // 3
			frontend.NewEOS(),
		},
		[][]int{
			{1, 3},
			{2, 3},
			{3},
			{4},
			{},
		})
}

func TestFollowMaybeNested3(x *testing.T) {
	testFollow(
		x,
		"q((bc?|x|y)?)*z",
		[]frontend.AST{
			frontend.NewCharacter('q'), // 0
			frontend.NewCharacter('b'), // 1
			frontend.NewCharacter('c'), // 2
			frontend.NewCharacter('x'), // 3
			frontend.NewCharacter('y'), // 4
			frontend.NewCharacter('z'), // 5
			frontend.NewEOS(),
		},
		[][]int{
			{1, 3, 4, 5},
			{2, 5, 1, 3, 4},
			{5, 1, 3, 4},
			{5, 1, 3, 4},
			{5, 1, 3, 4},
			{6},
			{},
		})
}

func TestFollowMaybeNested4(x *testing.T) {
	testFollow(
		x,
		"a(((bc?)?d?)?(ef?))?g",
		[]frontend.AST{
			frontend.NewCharacter('a'), // 0
			frontend.NewCharacter('b'), // 1
			frontend.NewCharacter('c'), // 2
			frontend.NewCharacter('d'), // 3
			frontend.NewCharacter('e'), // 4
			frontend.NewCharacter('f'), // 5
			frontend.NewCharacter('g'), // 6
			frontend.NewEOS(),
		},
		[][]int{
			{1, 3, 4, 6},
			{2, 3, 4},
			{3, 4},
			{4},
			{5, 6},
			{6},
			{7},
			{},
		})
}

func TestFollowNested(x *testing.T) {
	testFollow(
		x,
		"(a+b)*(c|d)*",
		[]frontend.AST{
			frontend.NewCharacter('a'), // 0
			frontend.NewCharacter('b'), // 1
			frontend.NewCharacter('c'), // 2
			frontend.NewCharacter('d'), // 3
			frontend.NewEOS(),
		},
		[][]int{
			{0, 1},
			{0, 2, 3, 4},
			{2, 3, 4},
			{2, 3, 4},
			{},
		})
}

func TestMatchesEmptyString_char(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a"))
	t.AssertNil(err)
	nullable := Label(ast).MatchesEmptyString()
	t.Assert(!nullable[len(nullable)-1], "character should not match the empty string, %v", ast)
}

func TestMatchesEmptyString_range(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("[a-z]"))
	t.AssertNil(err)
	nullable := Label(ast).MatchesEmptyString()
	t.Assert(!nullable[len(nullable)-1], "range should not match the empty string, %v", ast)
}

func TestMatchesEmptyString_maybe(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a?"))
	t.AssertNil(err)
	nullable := Label(ast).MatchesEmptyString()
	t.Assert(nullable[len(nullable)-1], "maybe should match the empty string, %v", ast)
}

func TestMatchesEmptyString_star(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a*"))
	t.AssertNil(err)
	nullable := Label(ast).MatchesEmptyString()
	t.Assert(nullable[len(nullable)-1], "star should match the empty string, %v", ast)
}

func TestMatchesEmptyString_plus(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a+"))
	t.AssertNil(err)
	nullable := Label(ast).MatchesEmptyString()
	t.Assert(!nullable[len(nullable)-1], "a+ should not match the empty string, %v", ast)
	ast, err = frontend.Parse([]byte("a?+"))
	t.AssertNil(err)
	nullable = Label(ast).MatchesEmptyString()
	t.Assert(nullable[len(nullable)-1], "a?+ should match the empty string, %v", ast)
}

func TestMatchesEmptyString_altMatch(x *testing.T) {
	t := (*test.T)(x)
	testMatchesEmptyString(t, "a|b", false, "a|b should not match the empty string, %v", "ab")
	testMatchesEmptyString(t, "a?|b", true, "a?|b should match the empty string, %v", "a?b")
	testMatchesEmptyString(t, "a|b?", true, "a|b? should match the empty string, %v", "ab?")
	testMatchesEmptyString(t, "a?|b?", true, "a?|b? should match the empty string, %v", "a?b?")
}

func TestMatchesEmptyString_concat(x *testing.T) {
	t := (*test.T)(x)
	testMatchesEmptyString(t, "ab", false, "ab should not match the empty string, %v", "ab")
	testMatchesEmptyString(t, "a?b", false, "a?b should not match the empty string, %v", "a?b")
	testMatchesEmptyString(t, "ab?", false, "ab? should not match the empty string, %v", "ab?")
	testMatchesEmptyString(t, "a?b?", true, "a?b? should match the empty string, %v", "a?b?")
}

func testMatchesEmptyString(t *test.T, regex string, matches bool, message string, args ...interface{}) {
	ast, err := frontend.Parse([]byte(regex))
	t.AssertNil(err)
	nullable := Label(ast).MatchesEmptyString()
	t.Assert(nullable[len(nullable)-1] == matches, message, args...)
}

func listEquals(a, b []frontend.AST) bool {
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

func astListOrder(lAst *LabeledAST, l []int) []frontend.AST {
	o := make([]frontend.AST, 0, len(l))
	for _, i := range l {
		o = append(o, lAst.Order[i])
	}
	return o
}

func astList(lAst *LabeledAST, l []int) []frontend.AST {
	o := make([]frontend.AST, 0, len(l))
	for _, i := range l {
		o = append(o, lAst.Order[lAst.Positions[i]])
	}
	return o
}

func TestFirst_char(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a"))
	t.AssertNil(err)
	lAst := Label(ast)
	first := []frontend.AST{
		frontend.NewCharacter('a'),
	}
	t.Assert(listEquals(first, astList(lAst, lAst.First()[len(lAst.Order)-1])), "first \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.First()[len(lAst.Order)-1]), first)
}

func TestLast_char(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a"))
	t.AssertNil(err)
	lAst := Label(ast)
	last := []frontend.AST{
		frontend.NewEOS(),
		frontend.NewCharacter('a'),
	}
	t.Assert(listEquals(last, astList(lAst, lAst.Last()[len(lAst.Order)-1])), "last \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.Last()[len(lAst.Order)-1]), last)
}

func TestFirst_range(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("[a-z]"))
	t.AssertNil(err)
	lAst := Label(ast)
	first := []frontend.AST{
		frontend.NewRange('a', 'z'),
	}
	t.Assert(listEquals(first, astList(lAst, lAst.First()[len(lAst.Order)-1])), "first \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.First()[len(lAst.Order)-1]), first)
}

func TestLast_range(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("[a-z]"))
	t.AssertNil(err)
	lAst := Label(ast)
	last := []frontend.AST{
		frontend.NewEOS(),
		frontend.NewRange('a', 'z'),
	}
	t.Assert(listEquals(last, astList(lAst, lAst.Last()[len(lAst.Order)-1])), "last \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.Last()[len(lAst.Order)-1]), last)
}

func TestFirst_ops(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a?*+"))
	t.AssertNil(err)
	lAst := Label(ast)
	first := []frontend.AST{
		frontend.NewCharacter('a'),
		frontend.NewEOS(),
	}
	t.Assert(listEquals(first, astList(lAst, lAst.First()[len(lAst.Order)-1])), "first \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.First()[len(lAst.Order)-1]), first)
}

func TestLast_ops(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a*?+"))
	t.AssertNil(err)
	lAst := Label(ast)
	last := []frontend.AST{
		frontend.NewEOS(),
		frontend.NewCharacter('a'),
	}
	t.Assert(listEquals(last, astList(lAst, lAst.Last()[len(lAst.Order)-1])), "last \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.Last()[len(lAst.Order)-1]), last)
}

func TestFirst_alt(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a|b"))
	t.AssertNil(err)
	lAst := Label(ast)
	first := []frontend.AST{
		frontend.NewCharacter('a'),
		frontend.NewCharacter('b'),
	}
	t.Assert(listEquals(first, astList(lAst, lAst.First()[len(lAst.Order)-1])), "first \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.First()[len(lAst.Order)-1]), first)
}

func TestLast_alt(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a|b"))
	t.AssertNil(err)
	lAst := Label(ast)
	last := []frontend.AST{
		frontend.NewEOS(),
		frontend.NewCharacter('a'),
		frontend.NewCharacter('b'),
	}
	t.Assert(listEquals(last, astList(lAst, lAst.Last()[len(lAst.Order)-1])), "last \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.Last()[len(lAst.Order)-1]), last)
}

func TestFirst_concat(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("a?b?c?de"))
	t.AssertNil(err)
	lAst := Label(ast)
	first := []frontend.AST{
		frontend.NewCharacter('a'),
		frontend.NewCharacter('b'),
		frontend.NewCharacter('c'),
		frontend.NewCharacter('d'),
	}
	t.Assert(listEquals(first, astList(lAst, lAst.First()[len(lAst.Order)-1])), "first \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.First()[len(lAst.Order)-1]), first)
}

func TestLast_concat(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("abc?d?e?"))
	t.AssertNil(err)
	lAst := Label(ast)
	last := []frontend.AST{
		frontend.NewEOS(),
		frontend.NewCharacter('e'),
		frontend.NewCharacter('d'),
		frontend.NewCharacter('c'),
		frontend.NewCharacter('b'),
	}
	t.Assert(listEquals(last, astList(lAst, lAst.Last()[len(lAst.Order)-1])), "last \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.Last()[len(lAst.Order)-1]), last)
}

func TestLast_concat2(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("abc*d*e*"))
	t.AssertNil(err)
	lAst := Label(ast)
	last := []frontend.AST{
		frontend.NewEOS(),
		frontend.NewCharacter('e'),
		frontend.NewCharacter('d'),
		frontend.NewCharacter('c'),
		frontend.NewCharacter('b'),
	}
	t.Assert(listEquals(last, astList(lAst, lAst.Last()[len(lAst.Order)-1])), "last \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.Last()[len(lAst.Order)-1]), last)
}

func TestLast_concat3(x *testing.T) {
	t := (*test.T)(x)
	ast, err := frontend.Parse([]byte("abc+d+e+"))
	t.AssertNil(err)
	lAst := Label(ast)
	last := []frontend.AST{
		frontend.NewEOS(),
		frontend.NewCharacter('e'),
	}
	t.Assert(listEquals(last, astList(lAst, lAst.Last()[len(lAst.Order)-1])), "last \n\tproduced: %v, \n\texpected: %v", astList(lAst, lAst.Last()[len(lAst.Order)-1]), last)
}
