package main

import (
	"fmt"
	"strings"

	"github.com/timtadh/lexmachine"
)

// Node is a very simple AST/Parse Tree node. It stores a Name (required), a
// Token (optional), and any child nodes.
type Node struct {
	Name     string
	Token    *lexmachine.Token
	Children []*Node
}

// NewNode makes a node from a name and a token. The token may be nil.
func NewNode(name string, token *lexmachine.Token) *Node {
	return &Node{
		Name:  name,
		Token: token,
	}
}

// AddKid puts a node at the end of the child list
func (n *Node) AddKid(kid *Node) *Node {
	n.Children = append(n.Children, kid)
	return n
}

// PrependKid puts a node at the beginning of the child list
func (n *Node) PrependKid(kid *Node) *Node {
	kids := append(make([]*Node, 0, cap(n.Children)+1), kid)
	n.Children = append(kids, n.Children...)
	return n
}

// String humanizes the tree starting at the current node.
func (n *Node) String() string {
	parts := make([]string, 0, len(n.Children))
	parts = append(parts, n.Name)
	if n.Token != nil && string(n.Token.Lexeme) != n.Name {
		parts = append(parts, fmt.Sprintf("%q", string(n.Token.Lexeme)))
	}
	for _, k := range n.Children {
		parts = append(parts, k.String())
	}
	if len(parts) > 1 {
		return fmt.Sprintf("(%v)", strings.Join(parts, " "))
	}
	return strings.Join(parts, " ")
}
