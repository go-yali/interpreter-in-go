/** A program in Monkey is a series of statements, which can be let statements, return statements, or expression statements **/

package ast

import (
	"bytes"
	"monkey/token"
)

// Each of our AST Nodes (ie each expression, statement) all must implement the Node interface, aka it must provide a TokenLiteral() that returns the literal value of the token it's associated with. TokenLiteral will only be used for debugging and testing
// All nodes will be connected to each other
// Some nodes implement the Statement interface, some the expression interface
type Node interface {
	TokenLiteral() string
	String() string
}

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

// LetStatement implements Statement interface
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct { // imple
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

// IntegerLiteral implements the Expression interface
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type PrefixExpression struct {
	Token    token.Token // a prefix token (! or -)
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    token.Token // The operator token, eg +
	Left     Expression
	Operator string
	Right    Expression
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// In order to add expression/let/return statements to the Statements slice of ast.Program, we satisfy the ast.Statement interface
func (ls *LetStatement) statementNode()        {}
func (rs *ReturnStatement) statementNode()     {}
func (i *Identifier) expressionNode()          {}
func (es *ExpressionStatement) statementNode() {}
func (il *IntegerLiteral) expressionNode()     {}
func (pe *PrefixExpression) expressionNode()   {}
func (ie *InfixExpression) expressionNode()    {}

func (ls *LetStatement) TokenLiteral() string        { return ls.Token.Literal }
func (i *Identifier) TokenLiteral() string           { return i.Token.Literal }
func (rs *ReturnStatement) TokenLiteral() string     { return rs.Token.Literal }
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (il *IntegerLiteral) TokenLiteral() string      { return il.Token.Literal }
func (pe *PrefixExpression) TokenLiteral() string    { return pe.Token.Literal }
func (ie *InfixExpression) TokenLiteral() string     { return ie.Token.Literal }

// Programs String method creates a buffer and writes the return value of each statement's String() method to it
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (i *Identifier) String() string {
	return i.Value
}

func (il *IntegerLiteral) String() string { return il.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	// Deliberately add paranthesses around the operator and Right,
	// allowing us to see which operand belongs to which operator
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
