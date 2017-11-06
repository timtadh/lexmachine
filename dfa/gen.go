package dfa

import (
	"fmt"
	"strings"

	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/linked"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
	"github.com/timtadh/lexmachine/frontend"
	"github.com/timtadh/lexmachine/machines"
)

type DFA struct {
	Start     int                   // the starting state
	Error     int                   // the error state (should be 0)
	Accepting machines.DFAAccepting // state-idx to match-id
	Trans     machines.DFATrans     // the transition matrix
	Matches   [][]int               // match-id to list of accepting states
}

// TODO
// 1. have Follow find the set of accepting positions grouped by match
// 2. add start states to DFA
// 3. add accepting states (grouped by match) to DFA
// 4. move dfa helpers here
// 5. make Follow more efficient
// 6. precompute first, last, and epsilon instead of using recursive defs
// 7. implement DFA minimization
// 8. then write a machine based on the DFA

func Generate(root frontend.AST) *DFA {
	ast := Label(root)
	positions := ast.Positions
	first, follow := ast.Follow()
	trans := hashtable.NewLinearHash()
	states := set.NewSortedSet(len(positions))
	accepting := set.NewSortedSet(len(positions))
	matchSet := set.NewSortedSet(len(ast.Matches))
	unmarked := linked.New()
	start := makeDState(first)
	trans.Put(start, make(map[byte]*set.SortedSet))
	states.Add(start)
	unmarked.Push(start)

	for _, m := range ast.Matches {
		matchSet.Add(types.Int(m))
	}

	for unmarked.Size() > 0 {
		x, err := unmarked.Pop()
		if err != nil {
			panic(err)
		}
		s := x.(*set.SortedSet)
		posBySymbol := make(map[int][]int)
		for pos, next := s.Items()(); next != nil; pos, next = next() {
			p := int(pos.(types.Int))
			switch n := ast.Order[positions[p]].(type) {
			case *frontend.EOS:
				posBySymbol[-1] = append(posBySymbol[-1], p)
			case *frontend.Character:
				sym := int(n.Char)
				posBySymbol[sym] = append(posBySymbol[sym], p)
			case *frontend.Range:
				for i := int(n.From); i <= int(n.To); i++ {
					posBySymbol[i] = append(posBySymbol[i], p)
				}
			}
		}
		for symbol, positions := range posBySymbol {
			if symbol == -1 {
				accepting.Add(s)
			} else if 0 <= symbol && symbol < 256 {
				// pFollow will be a new DState
				pFollow := set.NewSortedSet(len(positions) * 2)
				for _, p := range positions {
					for next := range follow[p] {
						pFollow.Add(types.Int(next))
					}
				}
				if !states.Has(pFollow) {
					trans.Put(pFollow, make(map[byte]*set.SortedSet))
					states.Add(pFollow)
					unmarked.Push(pFollow)
				}
				x, err := trans.Get(s)
				if err != nil {
					panic(err)
				}
				t := x.(map[byte]*set.SortedSet)
				t[byte(symbol)] = pFollow
			} else {
				panic("symbol outside of range")
			}
		}
	}

	idx := func(state *set.SortedSet) int {
		i, has, err := states.Find(state)
		if err != nil {
			panic(err)
		}
		if !has {
			panic(fmt.Errorf("missing state %v", state))
		}
		return i
	}

	dfa := &DFA{
		Start:     idx(start) + 1,
		Matches:   make([][]int, len(ast.Matches)),
		Accepting: make(machines.DFAAccepting),
		Trans:     make(machines.DFATrans, trans.Size()+1),
	}
	for k, v, next := trans.Iterate()(); next != nil; k, v, next = next() {
		from := k.(*set.SortedSet)
		toMap := v.(map[byte]*set.SortedSet)
		fromIdx := idx(from) + 1
		for symbol, to := range toMap {
			dfa.Trans[fromIdx][symbol] = idx(to) + 1
		}
		if accepting.Has(from) {
			idx := 0
			for ; idx < len(ast.Matches); idx++ {
				if from.Has(types.Int(ast.Matches[idx])) {
					break
				}
			}
			if idx >= len(ast.Matches) {
				panic(fmt.Errorf("Could not find any of %v in %v", ast.Matches, from))
			}
			dfa.Matches[idx] = append(dfa.Matches[idx], fromIdx)
			dfa.Accepting[fromIdx] = idx
		}
	}

	return dfa
}

func makeDState(positions []int) *set.SortedSet {
	s := set.NewSortedSet(len(positions))
	for _, p := range positions {
		s.Add(types.Int(p))
	}
	return s
}

func (dfa *DFA) String() string {
	lines := make([]string, 0, len(dfa.Trans))
	lines = append(lines, fmt.Sprintf("start: %d", dfa.Start))
	lines = append(lines, "accepting:")
	for i, matches := range dfa.Matches {
		t := make([]string, 0, len(matches))
		for _, m := range matches {
			t = append(t, fmt.Sprintf("%d", m))
		}
		lines = append(lines, fmt.Sprintf("    %d {%v}", i, strings.Join(t, ", ")))
	}
	for i, row := range dfa.Trans {
		t := make([]string, 0, 10)
		for sym, to := range row {
			if to == 0 {
				continue
			}
			t = append(t, fmt.Sprintf("%v->%v", string([]byte{byte(sym)}), to))
		}

		lines = append(lines, fmt.Sprintf("%d %v", i, strings.Join(t, ", ")))
	}
	return strings.Join(lines, "\n")
}

// match tests whether a text matches a single expression in the DFA. The
// match-id is returned on success, otherwise a -1 is the result.
func (dfa *DFA) match(text string) int {
	s := dfa.Start
	for tc := 0; tc < len(text); tc++ {
		s = dfa.Trans[s][text[tc]]
	}
	if mid, has := dfa.Accepting[s]; !has {
		return -1
	} else {
		return mid
	}
}
