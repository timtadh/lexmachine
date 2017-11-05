package dfa

import (
	"fmt"
	"strings"

	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/linked"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
	"github.com/timtadh/lexmachine/frontend"
)

type DFA struct {
	Start     int
	Accepting [][]int
	Trans     [][256]int
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

func Generate(ast frontend.AST) *DFA {
	lAst := Label(frontend.DesugarRanges(ast))
	positions := lAst.Positions
	first, follow := lAst.Follow()
	trans := hashtable.NewLinearHash()
	states := set.NewSortedSet(len(positions))
	accepting := set.NewSortedSet(len(positions))
	unmarked := linked.New()
	start := makeDState(first)
	trans.Put(start, make(map[byte]*set.SortedSet))
	states.Add(start)
	unmarked.Push(start)
	fmt.Println(follow)

	for unmarked.Size() > 0 {
		x, err := unmarked.Pop()
		if err != nil {
			panic(err)
		}
		s := x.(*set.SortedSet)
		posBySymbol := make(map[int][]int)
		for pos, next := s.Items()(); next != nil; pos, next = next() {
			p := int(pos.(types.Int))
			var sym int
			if char, is := lAst.Order[positions[p]].(*frontend.Character); is {
				sym = int(char.Char)
			} else if _, is := lAst.Order[positions[p]].(*frontend.EOS); is {
				sym = -1
			}
			posBySymbol[sym] = append(posBySymbol[sym], p)
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
		Start: idx(start) + 1,
		Trans: make([][256]int, trans.Size()+1),
	}
	fmt.Println("accepting", accepting)
	for k, v, next := trans.Iterate()(); next != nil; k, v, next = next() {
		from := k.(*set.SortedSet)
		toMap := v.(map[byte]*set.SortedSet)
		fmt.Println(from)
		fromIdx := idx(from) + 1
		for symbol, to := range toMap {
			fmt.Println("    ", symbol, to)
			dfa.Trans[fromIdx][symbol] = idx(to) + 1
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
