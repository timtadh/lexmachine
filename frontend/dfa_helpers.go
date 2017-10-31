package frontend

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
