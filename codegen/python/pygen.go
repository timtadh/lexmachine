package python

import (
	"fmt"
	"strings"

	"github.com/timtadh/lexmachine/dfa"
)

var header = `

def tokenize(input):
    return _Scanner(input).tokenize()

class Match(object):

    def __init__(self, match_id, lexeme):
        self.match_id = match_id
        self.lexeme = lexeme

    def __repr__(self):
        return self.__str__()

    def __str__(self):
        return "Match({}, {})".format(self.match_id, repr(self.lexeme))

class _Scanner(object):

    def __init__(self, input):
        self.input = input
        self.idx = 0
        self.buf = list()
        self.tokens = list()

    def tokenize(self):
        state = self.start()
        while state != None:
            state = state()
        if self.idx != len(self.input):
            self.eosError()
        return self.tokens

    def mvto(self, next_state):
        if self.idx >= len(self.input):
            raise Exception("internal DFA error, index out of bounds")
        self.buf.append(self.input[self.idx])
        self.idx += 1
        return next_state

    def match(self, match_id):
        self.tokens.append(Match(match_id, ''.join(self.buf)))
        self.buf = list()
        if self.idx < len(self.input):
            return self.start()
        return None

    def eosError(self, state=None):
        raise Exception("UnconsumedInput, {}".format(repr(self.input[self.idx-len(self.buf):])))

    def error(self, state, expected):
        raise Exception("UnexpectedInput, {}. expected one of: {}".format(
            repr(self.input[self.idx]), 
            [chr(x) for x in expected]))

`

func Generate(dfa *dfa.DFA) string {
	stateFuncs := make([]string, 0, len(dfa.Trans))
	stateFuncs = append(stateFuncs, genStart(dfa))
	for state := range dfa.Trans {
		stateFuncs = append(stateFuncs, genState(dfa, state))
	}
	return header + strings.Join(stateFuncs, "\n\n")
}

func genStart(dfa *dfa.DFA) string {
	lines := make([]string, 0, 3)
	lines = append(lines, fmt.Sprintf("    def start(self):"))
	lines = append(lines, fmt.Sprintf("        return self.state_%d", dfa.Start))
	return strings.Join(lines, "\n")
}

func genState(dfa *dfa.DFA, state int) string {
	trans := dfa.Trans[state]
	matchID, accepting := dfa.Accepting[state]
	lines := make([]string, 0, len(trans))
	lines = append(lines, fmt.Sprintf("    def state_%v(self):", state))
	if dfa.Error == state {
		lines = append(lines, fmt.Sprintf("        self.error(%d, [])", state))
		return strings.Join(lines, "\n")
	}
	if len(trans) > 0 && accepting {
		lines = append(lines, fmt.Sprintf("        if self.idx >= len(self.input):"))
		lines = append(lines, fmt.Sprintf("            return self.match(%v)", matchID))
	} else if len(trans) > 0 {
		lines = append(lines, fmt.Sprintf("        if self.idx >= len(self.input):"))
		lines = append(lines, fmt.Sprintf("            self.eosError(%v)", state))
		lines = append(lines, fmt.Sprintf("            return"))
	}
	first := true
	allowed := make([]string, 0, len(trans))
	for ord := 0; ord < len(trans); ord++ {
		if trans[ord] == dfa.Error {
			continue
		}
		allowed = append(allowed, fmt.Sprint(ord))
		if first {
			lines = append(lines, fmt.Sprintf("        if   ord(self.input[self.idx]) == %d:", ord))
			first = false
		} else {
			lines = append(lines, fmt.Sprintf("        elif ord(self.input[self.idx]) == %d:", ord))
		}
		lines = append(lines, fmt.Sprintf("            return self.mvto(self.state_%d)", trans[ord]))
	}
	if len(allowed) == 0 && !accepting && dfa.Error != state {
		panic("bad dfa")
	}
	if accepting {
		lines = append(lines, fmt.Sprintf("        return self.match(%d)", matchID))
	} else {
		lines = append(lines, fmt.Sprintf("        self.error(%d, [%v])", state, strings.Join(allowed, ", ")))
	}
	return strings.Join(lines, "\n")
}
