package frontend

import (
	"fmt"
	"strings"
)

type AST interface {
	String() string
}

type AltMatch struct {
	A AST
	B AST
}

func (self *AltMatch) String() string {
	return fmt.Sprintf("(AltMatch %v, %v)", self.A, self.B)
}

type Match struct {
	AST
}

func (self *Match) String() string {
	return fmt.Sprintf("(Match %v)", self.AST)
}

type Alternation struct {
	A AST
	B AST
}

func (self *Alternation) String() string {
	return fmt.Sprintf("(Alternation %v, %v)", self.A, self.B)
}

type Star struct {
	AST
}

func (self *Star) String() string {
	return fmt.Sprintf("(* %v)", self.AST)
}

type Plus struct {
	AST
}

func (self *Plus) String() string {
	return fmt.Sprintf("(+ %v)", self.AST)
}

type Maybe struct {
	AST
}

func (self *Maybe) String() string {
	return fmt.Sprintf("(? %v)", self.AST)
}

type Concat struct {
	Items []AST
}

func (self *Concat) String() string {
	s := "(Concat "
	items := make([]string, 0, len(self.Items))
	for _, i := range self.Items {
		items = append(items, i.String())
	}
	s += strings.Join(items, ", ") + ")"
	return s
}

type Range struct {
	From byte
	To   byte
}

func (self *Range) String() string {
	return fmt.Sprintf(
		"(Range %d %d)",
		self.From,
		self.To,
	)
}

type Character struct {
	Char byte
}

func (self *Character) String() string {
	return fmt.Sprintf(
		"(Character %s)",
		string([]byte{self.Char}),
	)
}

func NewAltMatch(a, b AST) AST {
	if a == nil || b == nil {
		panic("Alt match does not except nils")
	}
	return &AltMatch{a, b}
}

func NewMatch(ast AST) AST {
	return &Match{ast}
}

func NewAlternation(choice, alternation AST) AST {
	if alternation == nil {
		return choice
	}
	return &Alternation{choice, alternation}
}

func NewApplyOp(op, atomic AST) AST {
	switch o := op.(type) {
	case *Star:
		o.AST = atomic
	case *Plus:
		o.AST = atomic
	case *Maybe:
		o.AST = atomic
	default:
		panic("unexpected op")
	}
	return op
}

func NewOp(op string) AST {
	switch op {
	case "*":
		return &Star{}
	case "+":
		return &Plus{}
	case "?":
		return &Maybe{}
	default:
		panic("unexpected op")
	}
}

func NewConcat(char, concat AST) AST {
	if concat == nil {
		return char
	}
	if cc, ok := concat.(*Concat); ok {
		items := make([]AST, len(cc.Items)+1)
		items[0] = char
		for i, item := range cc.Items {
			items[i+1] = item
		}
		return &Concat{items}
	}
	return &Concat{[]AST{char, concat}}
}

func NewCharacter(b byte) AST {
	return &Character{b}
}

func NewAny() AST {
	return &Range{0, 255}
}

func NewRange(from, to AST) AST {
	return &Range{from.(*Character).Char, to.(*Character).Char}
}
