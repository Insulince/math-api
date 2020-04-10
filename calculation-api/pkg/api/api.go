package api

import (
	"math-api/calculation-api/pkg/configurations"
	"net/http"
	"math-api/base-api"
	"fmt"
	"log"
	"errors"
	"math-api/calculation-api/pkg/models"
)

type Api struct {
	Config *configurations.Config
}

func New(config *configurations.Config) (api *Api) {
	return &Api{Config: config}
}

func (a *Api) HomeHandler(ar *base_api.ApiRequest, aw *base_api.ApiResponseWriter) () {
	aw.Respond(struct{ Message string `json:"message"` }{Message: "Welcome!"}, http.StatusOK)
}

func (a *Api) HealthCheckHandler(ar *base_api.ApiRequest, aw *base_api.ApiResponseWriter) () {
	aw.Respond(struct{ Message string `json:"message"` }{Message: "OK"}, http.StatusOK)
}

func (a *Api) NotFoundHandler(ar *base_api.ApiRequest, aw *base_api.ApiResponseWriter) () {
	aw.Respond(struct{ Message string `json:"message"` }{Message: "Unsupported URL provided."}, http.StatusNotFound)
}

func (a *Api) CalculateHandler(ar *base_api.ApiRequest, aw *base_api.ApiResponseWriter) () {
	type Response struct {
		Message    string             `json:"message"`
		Expression *models.Expression `json:"expression"`
	}

	expressionQueryParameters := ar.GetQueryParameters()["expression"]
	if len(expressionQueryParameters) != 1 {
		err := errors.New(fmt.Sprintf("Could not parse expression: Invalid number of \"expression\" query parameters provided. Expected \"1\", but got \"%v\"", len(expressionQueryParameters)))
		log.Printf(err.Error())
		aw.Respond(&Response{Message: err.Error(), Expression: nil}, http.StatusBadRequest)
		return
	}
	expressionQueryParameter := expressionQueryParameters[0]

	expression := models.NewExpression(expressionQueryParameter, a.Config)

	err := expression.ParseExpressionString()
	if err != nil {
		log.Printf(err.Error())
		aw.Respond(&Response{Message: err.Error(), Expression: expression}, http.StatusBadRequest)
		return
	}
	fmt.Println(expression.ExpressionNodeRoot.ToString())

	err = expression.CalculateResult()
	if err != nil {
		log.Printf(err.Error())
		aw.Respond(&Response{Message: err.Error(), Expression: expression}, http.StatusInternalServerError)
		return
	}

	aw.Respond(&Response{Message: "Success.", Expression: expression}, http.StatusOK)
}
