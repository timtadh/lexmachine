package dfa

import (
	"testing"

	"github.com/timtadh/data-structures/test"
	"github.com/timtadh/lexmachine/frontend"
)

func mustParse(regex string) frontend.AST {
	ast, err := frontend.Parse([]byte(regex))
	if err != nil {
		panic(err)
	}
	return ast
}

func testGen(t *test.T, regex, text string, matchID int) {
	ast, err := frontend.Parse([]byte(regex))
	t.AssertNil(err)
	testGenMatch(t, ast, text, matchID)
}

func testGenMatch(t *test.T, ast frontend.AST, text string, matchID int) {
	dfa := Generate(ast)
	t.Assert(dfa.match(text) == matchID,
		"Expected match %d got %d for text %q.\nast: %v\n%v",
		matchID, dfa.match(text), text, ast, dfa)
}

func TestGenExample(x *testing.T) {
	t := (*test.T)(x)
	ast := mustParse("(([a-z]+[A-Z])*[0-9])?wizard")
	testGenMatch(t, ast, "abcAaAaZ0wizard", 0)
	testGenMatch(t, ast, "wizard", 0)
	testGenMatch(t, ast, "0wizard", 0)
	testGenMatch(t, ast, "7wizard", 0)
	testGenMatch(t, ast, "aaaA7wizard", 0)
	testGenMatch(t, ast, "A7wizard", -1)
	testGenMatch(t, ast, "a7wizard", -1)
	testGenMatch(t, ast, "Awizard", -1)
	testGenMatch(t, ast, "abcAaAaZ0wizar", -1)
}

func TestGenMin(x *testing.T) {
	// these regexes were tested/designed to give non-minimal dfas
	t := (*test.T)(x)
	testGen(t, "(a[a-c]*|a+)d", "abd", 0)
	testGen(t, "(a[a-c]*|a+)d", "accad", 0)
	testGen(t, "(a[a-c]*|a+)d", "a", -1)
	testGen(t, "(a[a-c]*|a+)d", "ad", 0)
	testGen(t, "(a[a-c]*|a+)d", "d", -1)

	ast := mustParse(`((((x[x-z]*c|x+c)+|x[u-z]*c|x+c)+|a[a-c]*c|a+c)+|a[a-c]*c|a+c)d`)
	testGenMatch(t, ast, "xcd", 0)
	testGenMatch(t, ast, "xzcd", 0)
	testGenMatch(t, ast, "acd", 0)
	testGenMatch(t, ast, "abcd", 0)
	testGenMatch(t, ast, "aaaacd", 0)
	testGenMatch(t, ast, "xxcd", 0)
	testGenMatch(t, ast, "xxzzzzcd", 0)
	testGenMatch(t, ast, "xxzuzzcd", 0)
	testGenMatch(t, ast, "xxzuzzaacd", -1)
}

func TestGenCharacter(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a", "a", 0)
	testGen(t, "a", "b", -1)
}

func TestGenRange(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "[a-c]", "a", 0)
	testGen(t, "[a-c]", "b", 0)
	testGen(t, "[a-c]", "c", 0)
	testGen(t, "[a-c]", "d", -1)
}

func TestGenMaybe(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a?", "a", 0)
	testGen(t, "a?", "", 0)
}

func TestGenMaybeNested1(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a(bc?)?d", "abcd", 0)
	testGen(t, "a(bc?)?d", "abd", 0)
	testGen(t, "a(bc?)?d", "ad", 0)
	testGen(t, "a(bc?)?d", "acd", -1)
	testGen(t, "a(bc?)?d", "abc", -1)
	testGen(t, "a(bc?)?d", "bc", -1)
	testGen(t, "a(bc?)?d", "bcd", -1)
}

func TestGenMaybeNested2(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a(b?c)?d", "abcd", 0)
	testGen(t, "a(b?c)?d", "acd", 0)
	testGen(t, "a(b?c)?d", "ad", 0)
	testGen(t, "a(b?c)?d", "a", -1)
	testGen(t, "a(b?c)?d", "d", -1)
	testGen(t, "a(b?c)?d", "ab", -1)
	testGen(t, "a(b?c)?d", "ac", -1)
	testGen(t, "a(b?c)?d", "abd", -1)
	testGen(t, "a(b?c)?d", "abd", -1)
}

func TestGenStar(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a*", "", 0)
	testGen(t, "a*", "aa", 0)
	testGen(t, "a*", "aaaaaaa", 0)
	testGen(t, "a*", "aaaabaaa", -1)
}

func TestGenStarNested1(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a(bc*)*d", "abcd", 0)
	testGen(t, "a(bc*)*d", "abd", 0)
	testGen(t, "a(bc*)*d", "ad", 0)
	testGen(t, "a(bc*)*d", "abcbbccccbbbbccd", 0)
	testGen(t, "a(bc*)*d", "acd", -1)
}

func TestGenStarNested2(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a(b*c)*d", "abcd", 0)
	testGen(t, "a(b*c)*d", "acd", 0)
	testGen(t, "a(b*c)*d", "ad", 0)
	testGen(t, "a(b*c)*d", "abcbbccccbbbbccd", 0)
	testGen(t, "a(b*c)*d", "abd", -1)
}

func TestGenPlus(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a+", "", -1)
	testGen(t, "a+", "aa", 0)
	testGen(t, "a+", "aaaaaaa", 0)
	testGen(t, "a+", "aaaabaaa", -1)
}

func TestGenAlt(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a|bcd|e", "a", 0)
	testGen(t, "a|bcd|e", "bcd", 0)
	testGen(t, "a|bcd|e", "e", 0)
	testGen(t, "a|bcd|e", "ae", -1)
	testGen(t, "a|bcd|e", "abcd", -1)
	testGen(t, "a|bcd|e", "bcde", -1)
	testGen(t, "a|bcd|e", "abcde", -1)
	testGen(t, "a|bcd|e", "a|bcd|e", -1)
	testGen(t, "a|bcd|e", "b", -1)
	testGen(t, "a|bcd|e", "c", -1)
	testGen(t, "a|bcd|e", "d", -1)
}

func TestGenAltAny(x *testing.T) {
	t := (*test.T)(x)
	testGen(t, "a|bcd|e|.", "a", 0)
	testGen(t, "a|bcd|e|.", "bcd", 0)
	testGen(t, "a|bcd|e|.", "e", 0)
	testGen(t, "a|bcd|e|.", "f", 0)
	testGen(t, "a|bcd|e|.", "\x00", 0)
	testGen(t, "a|bcd|e|.", "\x1f", 0)
	testGen(t, "a|bcd|e|.", "bcde", -1)
	testGen(t, "a|bcd|e|.", "b", 0)
	testGen(t, "a|bcd|e|.", "c", 0)
	testGen(t, "a|bcd|e|.", "d", 0)
}

func TestGenAltMatch(x *testing.T) {
	t := (*test.T)(x)
	ast := frontend.NewAltMatch(
		mustParse("b|e"),
		frontend.NewAltMatch(
			mustParse("[a-d]"),
			mustParse("f"),
		),
	)
	testGenMatch(t, ast, "a", 1)
	testGenMatch(t, ast, "b", 0)
	testGenMatch(t, ast, "c", 1)
	testGenMatch(t, ast, "d", 1)
	testGenMatch(t, ast, "e", 0)
	testGenMatch(t, ast, "f", 2)
	testGenMatch(t, ast, "A", -1)
}
