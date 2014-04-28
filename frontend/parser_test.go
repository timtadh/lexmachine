package frontend

import "testing"

func TestParse(t *testing.T) {
    ast, err := Parse([]byte("ab(a|c|d)?we*\\\\\\[\\..[s-f]+|qyx"))
    if err != nil {
        t.Error(err)
    }
    parsed := "(Alternation (Concat (Concat (Character a), (Character b)), (? (Alternation (Character a), (Alternation (Character c), (Character d)))), (* (Concat (Character w), (Character e))), (+ (Concat (Character \\), (Character [), (Character .), (Range 0 255), (Range 115 102)))), (Concat (Character q), (Character y), (Character x)))"
    if ast.String() != parsed {
        t.Log(ast.String())
        t.Log(parsed)
        t.Error("Did not parse correctly")
    }
}

