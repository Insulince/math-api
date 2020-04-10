package models

import (
	"errors"
	"regexp"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"math-api/base-api"
	"strconv"
	"math-api/calculation-api/pkg/configurations"
)

const ExpressionStringRegex = "^[0-9.\\+\\-\\/*()]+$"

type Expression struct {
	ExpressionString   string                 `json:"expressionString"`
	ExpressionNodeRoot *ExpressionNode        `json:"-"` /* Should not be JSON encoded */
	AsString           string                 `json:"asString"`
	CallsUsed          int                    `json:"callsUsed"`
	Result             float64                `json:"result"`
	Config             *configurations.Config `json:"-"` /* Should DEFINITELY not be JSON encoded */
}

func NewExpression(expressionString string, config *configurations.Config) (e *Expression) {
	return &Expression{ExpressionString: expressionString, ExpressionNodeRoot: nil, AsString: "", Config: config}
}

func (e *Expression) ParseExpressionString() (err error) {
	expressionConforms, err := regexp.MatchString(ExpressionStringRegex, e.ExpressionString)
	if err != nil {
		return err
	}
	if expressionConforms == false {
		return errors.New("Could not parse expression string, it did not match the expected format for an expression.\n")
	}

	rootExpressionNode, err := NewScanner(e.ExpressionString).Parse()
	if err != nil {
		return err
	}
	e.ExpressionNodeRoot = rootExpressionNode
	return nil
}

func (e *Expression) CalculateResult() (err error) {
	fmt.Println("API CALLS:")
	result, err := e.processExpressionNode(e.ExpressionNodeRoot)
	if err != nil {
		return err
	}
	e.Result = result
	e.AsString = fmt.Sprintf("%v = %v", e.ExpressionString, result)
	return nil
}

func (e *Expression) processExpressionNode(node *ExpressionNode) (result float64, err error) {
	if node.Type != T {
		if len(node.Children) == 1 {
			return e.processExpressionNode(node.Children[0])
		} else if len(node.Children) == 3 {
			switch node.Children[1].Value {
			case OperatorAdd:
				return e.calculateExpressionNode(node.Children[0], node.Children[2], "add")
			case OperatorSubtract:
				return e.calculateExpressionNode(node.Children[0], node.Children[2], "subtract")
			case OperatorMultiply:
				return e.calculateExpressionNode(node.Children[0], node.Children[2], "multiply")
			case OperatorDivide:
				return e.calculateExpressionNode(node.Children[0], node.Children[2], "divide")
			case DelimiterLeftParenthesis:
				return 0, errors.New(fmt.Sprintf("left paren: %v", node.Value))
			case DelimiterRightParenthesis:
				return 0, errors.New(fmt.Sprintf("right parent: %v", node.Value))
			case "":
				if node.Children[1].Type == T {
					return 0, errors.New(fmt.Sprintf("empty: %v", node.Value))
				}
				return e.processExpressionNode(node.Children[1])
			default:
				return 0, errors.New(fmt.Sprintf("_default %v", node.Value))
			}
		} else {
			return 0, errors.New(fmt.Sprintf("Unexpected number of children: %v", len(node.Children)))
		}
	} else {
		value, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
}

func (e *Expression) calculateExpressionNode(node1 *ExpressionNode, node2 *ExpressionNode, suffix string) (result float64, err error) {
	arg1, err := e.processExpressionNode(node2)
	if err != nil {
		return 0, err
	}
	arg2, err := e.processExpressionNode(node1)
	if err != nil {
		return 0, err
	}
	fmt.Printf("/%v/?arguments=%v,%v\n", suffix, arg1, arg2)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/%v/?arguments=%v,%v", e.Config.MathApiUrl, suffix, arg1, arg2), nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	rawResponseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var responseBody base_api.OperationResponse
	err = json.Unmarshal(rawResponseBody, &responseBody)
	if err != nil {
		return 0, err
	}

	e.CallsUsed++

	return responseBody.Result, nil
}
