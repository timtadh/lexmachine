package machines

import (
	"fmt"
	"sort"

	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/types"

	. "github.com/timtadh/lexmachine/inst"
	"github.com/timtadh/lexmachine/queue"
)

type dfa_state struct {
	id         int
	nfa_states pc_list
	moves
}

type dfa_stack []*dfa_state

func (stack dfa_stack) push(state *dfa_state) dfa_stack {
	return append(stack, state)
}

func (stack dfa_stack) pop() (dfa_stack, *dfa_state) {
	state := stack[len(stack)-1]
	return stack[:len(stack)-1], state
}

type pc_list []uint32

func (list pc_list) HasMatch(program InstSlice) bool {
	for _, pc := range list {
		if program[pc].Op == MATCH {
			return true
		}
	}
	return false
}

func (list pc_list) EarliestMatch(program InstSlice) uint32 {
	match := uint32(len(program) + 1)
	for _, pc := range list {
		if program[pc].Op == MATCH {
			if pc < match {
				match = pc
			}
		}
	}
	return match
}

func (list pc_list) Equals(o types.Equatable) bool {
	if l, ok := o.(pc_list); ok {
		if len(list) != len(l) {
			return false
		}
		for i := 0; i < len(list); i++ {
			if list[i] != l[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (list pc_list) Less(o types.Sortable) bool {
	if l, ok := o.(pc_list); ok {
		if len(list) < len(l) {
			return true
		} else if len(list) > len(l) {
			return false
		}
		for i := 0; i < len(list); i++ {
			if list[i] >= l[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (list pc_list) Hash() int {
	hash := 0
	for i, pc := range list {
		hash += (i + 1) * (int(pc) + 1)
	}
	return hash
}

func (list pc_list) has(pc uint32) bool {
	var l int = 0
	var r int = len(list) - 1
	var m int
	for l <= r {
		m = ((r - l) >> 1) + l
		if pc < list[m] {
			r = m - 1
		} else if pc == list[m] {
			return true
		} else {
			l = m + 1
		}
	}
	return false
}

func (list pc_list) insert(pc uint32) pc_list {
	for i := 0; i < len(list); i += 1 {
		if list[i] == pc {
			return list
		} else if pc < list[i] {
			return list.insert_at(i, pc)
		}
	}
	return list.insert_at(len(list), pc)
}

func (list pc_list) insert_at(i int, pc uint32) pc_list {
	var c uint32 = pc
	for ; i < len(list); i += 1 {
		c, list[i] = list[i], c
	}
	return append(list, c)
}

func closure_one(program InstSlice, pc uint32) pc_list {
	return closure(program, pc_list([]uint32{pc}))
}

func closure(program InstSlice, set pc_list) pc_list {
	list := make(pc_list, 0, 10)
	q := queue.New(len(program))
	for _, pc := range set {
		q.Push(pc)
	}
	for !q.Empty() {
		pc := q.Pop()
		list = list.insert(pc)
		inst := program[pc]
		switch inst.Op {
		case CHAR: // no actions are further reachable
		case MATCH: // no actions are further reachable
		case SPLIT:
			if !list.has(inst.Y) {
				q.Push(inst.Y)
			}
			fallthrough
		case JMP:
			if !list.has(inst.X) {
				q.Push(inst.X)
			}
		}
	}
	return list
}

type movement struct {
	a, b byte
	U    pc_list
}

func (m *movement) String() string {
	return fmt.Sprintf("<%v %v %v>", string([]byte{m.a}), string([]byte{m.b}), m.U)
}

type moves []*movement

func move(program InstSlice, T pc_list) (m moves) {
	for _, pc := range T {
		inst := program[pc]
		if inst.Op == CHAR {
			m = append(m, &movement{
				a: byte(inst.X),
				b: byte(inst.Y),
				U: closure_one(program, pc+1),
			})
		} else {
			// no other operation has movements!
		}
	}
	sort.Slice(m, func(i, j int) bool {
		min_i := m[i].U[0]
		min_j := m[j].U[0]
		return min_i < min_j
	})
	return m
}

func ToDFA(program InstSlice) InstSlice {
	dfa_states := hashtable.NewLinearHash()
	stack := make(dfa_stack, 0, 10)

	next_id := 0
	s0 := &dfa_state{id: next_id, nfa_states: closure_one(program, 0)}
	stack = stack.push(s0)
	if err := dfa_states.Put(s0.nfa_states, s0); err != nil {
		panic(err)
	}
	next_id++

	for len(stack) > 0 {
		var S *dfa_state
		stack, S = stack.pop()
		S.moves = move(program, S.nfa_states)
		for _, M := range S.moves {
			if !dfa_states.Has(M.U) {
				s := &dfa_state{id: next_id, nfa_states: M.U}
				next_id++
				if err := dfa_states.Put(M.U, s); err != nil {
					panic(err)
				}
				stack = stack.push(s)
			}
		}
	}

	dfa_keys := make([]pc_list, 0, 10)
	for k, next := dfa_states.Keys()(); next != nil; k, next = next() {
		dfa_keys = append(dfa_keys, k.(pc_list))
	}
	sort.Slice(dfa_keys, func(i, j int) bool {
		min_i := dfa_keys[i][0]
		min_j := dfa_keys[j][0]
		return min_i < min_j
	})
	fmt.Println(dfa_keys)

	dfa_build := make([]InstSlice, dfa_states.Size()+1)
	for _, k := range dfa_keys {
		v, _ := dfa_states.Get(k)
		s := v.(*dfa_state)
		fmt.Println("-", s.moves)
		for _, move := range s.moves {
			fmt.Println("move", move)
			u, err := dfa_states.Get(move.U)
			if err != nil {
				panic(err)
			}
			uid := uint32(u.(*dfa_state).id)
			fmt.Println("dfa state", uid)
			dfa_build[s.id] = append(dfa_build[s.id], &Inst{CHJMP, uint32(move.a), uint32(move.b)})
			dfa_build[s.id] = append(dfa_build[s.id], &Inst{JMP, uid, 0})
		}
		// TODO: track the NFA state the MATCH jump is coming from. The Lexer
		// engine needs this to communicate which pattern the MATCH corresponds
		// to.
		if s.nfa_states.HasMatch(program) {
			dfa_build[s.id] = append(dfa_build[s.id], &Inst{MATCH, s.nfa_states.EarliestMatch(program), 0})
			// dfa_build[s.id] = append(dfa_build[s.id], &Inst{JMP, uint32(len(dfa_build)-1), 0})
		}

		/*
			fmt.Println(k, s.id)
			for _, inst := range dfa_build[s.id] {
				fmt.Println("    ", inst)
			} */
	}
	// dfa_build[len(dfa_build)-1] = append(dfa_build[len(dfa_build)-1], &Inst{MATCH, 0, 0})

	dfa := make(InstSlice, 0, len(program))
	dfajmp := make([]int, len(dfa_build))
	for i, insts := range dfa_build {
		dfajmp[i] = len(dfa)
		for _, inst := range insts {
			dfa = append(dfa, inst)
		}
	}

	for _, inst := range dfa {
		if inst.Op == JMP {
			inst.X = uint32(dfajmp[inst.X])
		}
	}
	return dfa
}

func remove_unneeded_chjmp(dfa InstSlice) InstSlice {
	count := func(i uint32, removed []uint32) uint32 {
		count := uint32(0)
		for _, x := range removed {
			if x < i {
				count++
			}
		}
		return count
	}
	fmt.Println(dfa)
	optimized := make(InstSlice, 0, len(dfa))
	removed := make([]uint32, 0, 10)
	for i := 0; i < len(dfa); i++ {
		if dfa[i].Op == CHJMP && i+2 < len(dfa) {
			if dfa[i+1].Op == JMP && dfa[i+1].X == uint32(i+2) {
				char := &Inst{CHAR, dfa[i].X, dfa[i].Y}
				optimized = append(optimized, char)
				fmt.Println("Useless Jump", i, dfa[i], char)
				i++
				removed = append(removed, uint32(i+1))
				continue
			}
		}
		optimized = append(optimized, dfa[i])
	}
	for _, i := range optimized {
		switch i.Op {
		case SPLIT:
			i.Y = i.Y - count(i.Y, removed)
			fallthrough
		case JMP:
			i.X = i.X - count(i.X, removed)
		}
	}
	fmt.Println("removed", removed)
	return optimized
}
