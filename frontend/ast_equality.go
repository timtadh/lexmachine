package frontend

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

func (e *EOS) Equals(o AST) bool {
	if x, is := o.(*EOS); is {
		return *e == *x
	} else {
		return false
	}
}
