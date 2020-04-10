package models

import "fmt"

const RootLevel = 0

const /* NodeType */ (
	S = "S"
	E = "E"
	A = "A"
	M = "M"
	P = "P"
	N = "N"
	T = "T"
)

type ExpressionNode struct {
	Type       string /* NodeType */ `json:"type"`
	IsTerminal bool                  `json:"isTerminal"`
	Level      int                   `json:"level"`
	Value      string                `json:"value"`
	Children   []*ExpressionNode     `json:"children"`
}

func NewExpressionNode(tipe string, level int) (newExpressionNode *ExpressionNode) {
	if tipe != T {
		return &ExpressionNode{Type: tipe, IsTerminal: false, Level: level, Value: "", Children: make([]*ExpressionNode, 0)}
	} else {
		return &ExpressionNode{Type: tipe, IsTerminal: true, Level: level, Value: "", Children: make([]*ExpressionNode, 0)}
	}
}

func (en *ExpressionNode) ToString() (output string) {
	output += fmt.Sprintf("PARSE TREE:\n%v | %v\n", en.Type, en.Value)
	for _, child := range en.Children {
		output += child.toString(1)
	}
	return output
}

func (en *ExpressionNode) toString(level int) (output string) {
	spacer := ""
	for i := 0; i < level; i++ {
		spacer += "\t"
	}
	output += fmt.Sprintf("%v%v | %v\n", spacer, en.Type, en.Value)
	for _, child := range en.Children {
		output += child.toString(level + 1)
	}
	return output
}
