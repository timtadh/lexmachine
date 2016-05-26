package frontend

import (
	"fmt"
)

import (
	"github.com/timtadh/lexmachine/inst"
)

type generator struct {
	program inst.InstSlice
}

func Generate(ast AST) (inst.InstSlice, error) {
	g := &generator{
		program: make([]*inst.Inst, 0, 100),
	}
	fill := g.gen(ast)
	if len(fill) != 0 {
		return nil, fmt.Errorf("unconnected instructions")
	}
	return g.program, nil
}

func (self *generator) gen(ast AST) (fill []*uint32) {
	switch n := ast.(type) {
	case *AltMatch:
		fill = self.altMatch(n)
	case *Match:
		fill = self.match(n)
	case *Alternation:
		fill = self.alt(n)
	case *Star:
		fill = self.star(n)
	case *Plus:
		fill = self.plus(n)
	case *Maybe:
		fill = self.maybe(n)
	case *Concat:
		fill = self.concat(n)
	case *Character:
		fill = self.character(n)
	case *Range:
		fill = self.rangeGen(n)
	}
	return fill
}

func (self *generator) dofill(fill []*uint32) {
	for _, jmp := range fill {
		*jmp = uint32(len(self.program))
	}
}

func (self *generator) altMatch(a *AltMatch) []*uint32 {
	split := inst.New(inst.SPLIT, 0, 0)
	self.program = append(self.program, split)
	split.X = uint32(len(self.program))
	self.gen(a.A)
	split.Y = uint32(len(self.program))
	self.gen(a.B)
	return nil
}

func (self *generator) match(m *Match) []*uint32 {
	self.dofill(self.gen(m.AST))
	self.program = append(
		self.program, inst.New(inst.MATCH, 0, 0))
	return nil
}

func (self *generator) alt(a *Alternation) (fill []*uint32) {
	split := inst.New(inst.SPLIT, 0, 0)
	self.program = append(self.program, split)
	split.X = uint32(len(self.program))
	self.dofill(self.gen(a.A))
	jmp := inst.New(inst.JMP, 0, 0)
	self.program = append(self.program, jmp)
	split.Y = uint32(len(self.program))
	fill = self.gen(a.B)
	fill = append(fill, &jmp.X)
	return fill
}

func (self *generator) repeat(ast AST) (fill []*uint32) {
	split := inst.New(inst.SPLIT, 0, 0)
	split_pos := uint32(len(self.program))
	self.program = append(self.program, split)
	split.X = uint32(len(self.program))
	self.dofill(self.gen(ast))
	jmp := inst.New(inst.JMP, split_pos, 0)
	self.program = append(self.program, jmp)
	return []*uint32{&split.Y}
}

func (self *generator) star(s *Star) (fill []*uint32) {
	return self.repeat(s.AST)
}

func (self *generator) plus(p *Plus) (fill []*uint32) {
	self.dofill(self.gen(p.AST))
	return self.repeat(p.AST)
}

func (self *generator) maybe(m *Maybe) (fill []*uint32) {
	split := inst.New(inst.SPLIT, 0, 0)
	self.program = append(self.program, split)
	split.X = uint32(len(self.program))
	fill = self.gen(m.AST)
	fill = append(fill, &split.Y)
	return fill
}

func (self *generator) concat(c *Concat) (fill []*uint32) {
	for _, ast := range c.Items {
		self.dofill(fill)
		fill = self.gen(ast)
	}
	return fill
}

func (self *generator) character(ch *Character) []*uint32 {
	self.program = append(
		self.program,
		inst.New(inst.CHAR, uint32(ch.Char), uint32(ch.Char)))
	return nil
}

func (self *generator) rangeGen(r *Range) []*uint32 {
	self.program = append(
		self.program,
		inst.New(inst.CHAR, uint32(r.From), uint32(r.To)))
	return nil
}
