/*
If you ever get the INSANE idea to try to refactor this scanner into a left-to-right reading left-associative recursive decent parser,
please revisit this link and do us both a favor and DON'T attempt that:
https://www.gamedev.net/forums/topic/416784-recursive-descent-parsing-handling-left-associativity/

TOKENS
"0", "1", "2", "3", "4",
"5", "6", "7", "8", "9",
"(", ")", "+", "-", "*",
"/"

BNF
0 N <S> -> <E> | <A>
1 T <E> ->
2 N <A> -> <M> + <A> | <M> - <A> | <M>
3 N <M> -> <P> * <M> | <P> / <M> | <P>
4 N <P> -> ( <A> ) | <P>
5 T <N> -> [0-9]+
*/

package models

import (
	"errors"
	"fmt"
	"math-api/calculation-api/pkg/util"
	"strconv"
)

const /* Tokens */ (
	OperatorAdd               = "+"
	OperatorSubtract          = "-"
	OperatorMultiply          = "*"
	OperatorDivide            = "/"
	DelimiterLeftParenthesis  = "("
	DelimiterRightParenthesis = ")"
	/* And Integers */
)

type Scanner struct {
	Input        string `json:"input"`
	Position     int    `json:"position"`
	CurrentToken string `json:"currentToken"`
}

func NewScanner(input string) (s *Scanner) {
	s = &Scanner{Input: input, Position: len(input) - 1}
	s.nextToken()
	return s
}

func (s *Scanner) currentCharacter() (character string) {
	return string(s.Input[s.Position])
}

func (s *Scanner) advancePosition() () {
	s.Position--
}

func (s *Scanner) rollbackPosition() () {
	s.Position++
}

func (s *Scanner) consumeCharacter() (character string) {
	character = s.currentCharacter()
	s.advancePosition()
	return character
}

func (s *Scanner) moreCharactersRemain() (currentCharacterIsConsumable bool) {
	return s.Position >= 0
}

func (s *Scanner) nextToken() (err error) {
	mirrorToken := ""
	token := ""
	previousCharacter := ""
	tokenIndeterminate := true

	for s.moreCharactersRemain() && tokenIndeterminate {
		currentCharacter := s.consumeCharacter()

		if previousCharacter == "" {
			switch currentCharacter {
			case OperatorAdd:
				fallthrough
			case OperatorSubtract:
				fallthrough
			case OperatorMultiply:
				fallthrough
			case OperatorDivide:
				fallthrough
			case DelimiterLeftParenthesis:
				fallthrough
			case DelimiterRightParenthesis:
				tokenIndeterminate = false
			default:
				if characterIsNumeric(currentCharacter) {
					previousCharacter = currentCharacter
				} else {
					return errors.New(fmt.Sprintf("Unrecognized character encountered while building token: \"%v\"", currentCharacter))
				}
			}
			mirrorToken = currentCharacter
		} else {
			if characterIsNumeric(currentCharacter) {
				previousCharacter = currentCharacter
				mirrorToken += currentCharacter
			} else {
				tokenIndeterminate = false
				s.rollbackPosition()
			}
		}
	}

	for i := len(mirrorToken) - 1; i >= 0; i-- {
		token += string(mirrorToken[i])
	}
	s.CurrentToken = token
	return nil
}

func characterIsNumeric(character string) (isNumeric bool) {
	_, err := strconv.ParseInt(character, 10, 64)
	if err != nil {
		if character == "." {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func (s *Scanner) peek() (token string) {
	clone := &Scanner{Input: s.Input, Position: s.Position, CurrentToken: s.CurrentToken}
	clone.nextToken()
	return clone.CurrentToken
}

func (s *Scanner) consumeNonTerminal(node *ExpressionNode, child *ExpressionNode) () {
	if child.IsTerminal == false {
		node.Children = append(node.Children, child)
	}
}

func (s *Scanner) consumeTerminal(node *ExpressionNode) (err error) {
	child := NewExpressionNode(T, node.Level+1)
	child.Value = s.CurrentToken
	node.Children = append(node.Children, child)

	err = s.nextToken()
	if err != nil {
		return err
	}
	return nil
}

func (s *Scanner) Parse() (root *ExpressionNode, err error) {
	root = NewExpressionNode(S, RootLevel)
	if s.CurrentToken != "" {
		addSubNode, err := s.parseAddSub(RootLevel + 1)
		if err != nil {
			return nil, err
		}
		s.consumeNonTerminal(root, addSubNode)
	} else {
		err = s.consumeTerminal(root)
		if err != nil {
			return nil, err
		}
	}

	return root, nil
}

func (s *Scanner) parseAddSub(level int) (node *ExpressionNode, err error) {
	node = NewExpressionNode(A, level)
	miltDivNode, err := s.parseMultDiv(level + 1)
	if err != nil {
		return nil, err
	}
	s.consumeNonTerminal(node, miltDivNode)
	if s.CurrentToken == OperatorAdd || s.CurrentToken == OperatorSubtract {
		err = s.consumeTerminal(node)
		if err != nil {
			return nil, err
		}
		addSubNode, err := s.parseAddSub(level + 1)
		if err != nil {
			return nil, err
		}
		s.consumeNonTerminal(node, addSubNode)
	}
	return node, nil
}

func (s *Scanner) parseMultDiv(level int) (node *ExpressionNode, err error) {
	node = NewExpressionNode(M, level)
	parenNode, err := s.parseParen(level + 1)
	if err != nil {
		return nil, err
	}
	s.consumeNonTerminal(node, parenNode)
	if s.CurrentToken == OperatorMultiply || s.CurrentToken == OperatorDivide {
		err = s.consumeTerminal(node)
		if err != nil {
			return nil, err
		}
		multDivNode, err := s.parseMultDiv(level + 1)
		if err != nil {
			return nil, err
		}
		s.consumeNonTerminal(node, multDivNode)
	}
	return node, nil
}

func (s *Scanner) parseParen(level int) (node *ExpressionNode, err error) {
	node = NewExpressionNode(P, level)
	if s.CurrentToken == DelimiterRightParenthesis {
		err = s.consumeTerminal(node)
		if err != nil {
			return nil, err
		}
		addSubNode, err := s.parseAddSub(level + 1)
		if err != nil {
			return nil, err
		}
		s.consumeNonTerminal(node, addSubNode)
		if s.CurrentToken != DelimiterLeftParenthesis {
			return nil, errors.New(fmt.Sprintf("Missing closing parenthesis: %v", s.CurrentToken))
		}
		err = s.consumeTerminal(node)
		if err != nil {
			return nil, err
		}
	} else if util.IsNumber(s.CurrentToken) {
		err = s.consumeTerminal(node)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Unrecognized token: %v", s.CurrentToken))
	}
	return node, nil
}
