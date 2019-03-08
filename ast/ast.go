/** A program in Monkey is a series of statements, which can be let statements, return statements, or expression statements **/

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

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token       token.Token // the 'return token'
	ReturnValue Expression
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// In order to add expression/let/return statements to the Statements slice of ast.Program, we satisfy the ast.Statement interface
func (ls *LetStatement) statementNode()       {}
func (rs *ReturnStatement) statementNode()    {}
func (i *Identifier) expressionNode()         {}
func (es *ExpressionStatement) statementNoe() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
