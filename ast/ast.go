package ast

import "monkey/token"

// Each of our AST Nodes (ie each expression, statement) all must implement the Node interface, aka it must provide a TokenLiteral() that returns the literal value of the token it's associated with. TokenLiteral will only be used for debugging and testing
// All nodes will be connected to each other
// Some nodes implement the Statement interface, some the expression interface
type Node interface {
	TokenLiteral() string
}

// Statement interface only contains a dummy method, statementNode()
type Statement interface {
	Node
	statementNode()
}

// Expression interface only contains a dummy method, expressionNode()
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (ls *LetStatement) statementNode() {}
func (i *Identifier) expressionNode()   {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
