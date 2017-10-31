package frontend

import (
	"fmt"
	"unsafe"
)

func Follow(root AST) (positions []AST, follow []map[int]bool) {
	positions = findPositions(root)
	posmap := make(map[uintptr]int)
	ptr := func(p AST) uintptr {
		switch n := p.(type) {
		case *Range:
			return uintptr(unsafe.Pointer(n))
		case *Character:
			return uintptr(unsafe.Pointer(n))
		default:
			panic(fmt.Errorf("%T is not a Range or Character", p))
		}
	}
	for i, p := range positions {
		posmap[ptr(p)] = i
	}
	pos := func(p AST) int {
		return posmap[ptr(p)]
	}
	follow = make([]map[int]bool, len(positions))
	for i := range follow {
		follow[i] = make(map[int]bool)
	}
	stack := make([]AST, 0, 10)
	stack = append(stack, root)
	maybes := make([]*Maybe, 0, 10)
	for len(stack) > 0 {
		var cur AST
		stack, cur = stack[:len(stack)-1], stack[len(stack)-1]
		stack = append(stack, cur.Children()...)
		switch n := cur.(type) {
		case *Maybe:
			maybes = append(maybes, n)
		case *Concat:
			for x := 0; x < len(n.Items)-1; x++ {
				a := n.Items[x]
				b := n.Items[x+1]
				bFirst := make([]int, 0, 10)
				for _, p := range b.First() {
					bFirst = append(bFirst, pos(p))
				}
				for _, q := range a.Last() {
					i := pos(q)
					for _, j := range bFirst {
						follow[i][j] = true
					}
				}
			}
		case *Star, *Plus:
			nFirst := make([]int, 0, 10)
			for _, p := range n.First() {
				nFirst = append(nFirst, pos(p))
			}
			for _, q := range n.Last() {
				i := pos(q)
				for _, j := range nFirst {
					follow[i][j] = true
				}
			}
		}
	}

	subtree := func(n AST) map[int]bool {
		tree := make(map[int]bool)
		for _, p := range findPositions(n) {
			tree[pos(p)] = true
		}
		return tree
	}

	// Fix the maybes
	// everything with First(Maybe) in their FOLLOW
	//   (not including things in First(Maybe) in the Maybe subtree).
	// also has FOLLOW(Last(Maybe)) in their FOLLOW
	//   (not including things in Follow(Last(Maybe)) in the Maybe subtree).
	for i := 0; i < len(maybes); i++ {
		maybe := maybes[i]
		tree := subtree(maybe)
		for _, first := range maybe.First() {
			f := pos(first)
			// fmt.Printf("first%v = %v %v\n", maybe, first, f)
			for j, row := range follow {
				if _, has := tree[j]; has {
					// ignore things inside the Maybe subtree
					continue
				}
				if _, has := row[f]; has {
					for _, last := range maybe.Last() {
						l := pos(last)
						for following := range follow[l] {
							if _, has := tree[following]; has {
								// ignore things inside the Maybe subtree
								continue
							}
							follow[j][following] = true
							// fmt.Printf("%v -> %v\n", j, following)
						}
					}
				}
			}
		}
	}

	return positions, follow
}

func findPositions(ast AST) []AST {
	positions := make([]AST, 0, 10)
	stack := make([]AST, 0, 10)
	stack = append(stack, ast)
	for len(stack) > 0 {
		var cur AST
		stack, cur = stack[:len(stack)-1], stack[len(stack)-1]
		kids := cur.Children()
		for i := len(kids) - 1; i >= 0; i-- {
			stack = append(stack, kids[i])
		}
		switch cur.(type) {
		case *Character:
			positions = append(positions, cur)
		case *Range:
			positions = append(positions, cur)
		}
	}
	return positions
}

func (a *AltMatch) MatchesEmptyString() bool {
	return a.A.MatchesEmptyString() || a.B.MatchesEmptyString()
}

func (a *Alternation) MatchesEmptyString() bool {
	return a.A.MatchesEmptyString() || a.B.MatchesEmptyString()
}

func (c *Concat) MatchesEmptyString() bool {
	for _, i := range c.Items {
		if !i.MatchesEmptyString() {
			return false
		}
	}
	return true
}

func (m *Match) MatchesEmptyString() bool     { return m.AST.MatchesEmptyString() }
func (s *Star) MatchesEmptyString() bool      { return true }
func (p *Plus) MatchesEmptyString() bool      { return p.AST.MatchesEmptyString() }
func (m *Maybe) MatchesEmptyString() bool     { return true }
func (c *Character) MatchesEmptyString() bool { return false }
func (r *Range) MatchesEmptyString() bool     { return false }

func (c *Concat) First() []AST {
	first := make([]AST, 0, len(c.Items))
	for _, item := range c.Items {
		first = append(first, item.First()...)
		if !item.MatchesEmptyString() {
			break
		}
	}
	return first
}

func (a *AltMatch) First() []AST    { return append(a.A.First(), a.B.First()...) }
func (a *Alternation) First() []AST { return append(a.A.First(), a.B.First()...) }
func (m *Match) First() []AST       { return m.AST.First() }
func (s *Star) First() []AST        { return s.AST.First() }
func (p *Plus) First() []AST        { return p.AST.First() }
func (m *Maybe) First() []AST       { return m.AST.First() }
func (c *Character) First() []AST   { return []AST{c} }
func (r *Range) First() []AST       { return []AST{r} }

func (c *Concat) Last() []AST {
	last := make([]AST, 0, len(c.Items))
	for i := len(c.Items) - 1; i >= 0; i-- {
		item := c.Items[i]
		last = append(last, item.Last()...)
		if !item.MatchesEmptyString() {
			break
		}
	}
	return last
}

func (a *AltMatch) Last() []AST    { return append(a.A.Last(), a.B.Last()...) }
func (a *Alternation) Last() []AST { return append(a.A.Last(), a.B.Last()...) }
func (m *Match) Last() []AST       { return m.AST.Last() }
func (s *Star) Last() []AST        { return s.AST.Last() }
func (p *Plus) Last() []AST        { return p.AST.Last() }
func (m *Maybe) Last() []AST       { return m.AST.Last() }
func (c *Character) Last() []AST   { return []AST{c} }
func (r *Range) Last() []AST       { return []AST{r} }

func (c *Concat) Equals(o AST) bool {
	if x, is := o.(*Concat); is {
		if len(c.Items) != len(x.Items) {
			return false
		}
		for i := range c.Items {
			if !c.Items[i].Equals(x.Items[i]) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (a *AltMatch) Equals(o AST) bool {
	if x, is := o.(*AltMatch); is {
		return a.A.Equals(x.A) && a.B.Equals(x.B)
	} else {
		return false
	}
}

func (a *Alternation) Equals(o AST) bool {
	if x, is := o.(*Alternation); is {
		return a.A.Equals(x.A) && a.B.Equals(x.B)
	} else {
		return false
	}
}

func (m *Match) Equals(o AST) bool {
	if x, is := o.(*Match); is {
		return m.AST.Equals(x.AST)
	} else {
		return false
	}
}

func (s *Star) Equals(o AST) bool {
	if x, is := o.(*Star); is {
		return s.AST.Equals(x.AST)
	} else {
		return false
	}
}

func (p *Plus) Equals(o AST) bool {
	if x, is := o.(*Plus); is {
		return p.AST.Equals(x.AST)
	} else {
		return false
	}
}

func (m *Maybe) Equals(o AST) bool {
	if x, is := o.(*Maybe); is {
		return m.AST.Equals(x.AST)
	} else {
		return false
	}
}

func (c *Character) Equals(o AST) bool {
	if x, is := o.(*Character); is {
		return *c == *x
	} else {
		return false
	}
}

func (r *Range) Equals(o AST) bool {
	if x, is := o.(*Range); is {
		return *r == *x
	} else {
		return false
	}
}
