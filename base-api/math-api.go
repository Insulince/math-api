package base_api

import (
	"net/http"
	"log"
	"fmt"
	"errors"
	"strconv"
	"strings"
)

type Api struct {
	Config *Config
}

func New(config *Config) (api *Api) {
	return &Api{Config: config}
}

func (a *Api) HomeHandler(ar *ApiRequest, aw *ApiResponseWriter) () {
	aw.Respond(struct{ Message string `json:"message"` }{Message: "Welcome!"}, http.StatusOK)
}

func (a *Api) HealthCheckHandler(ar *ApiRequest, aw *ApiResponseWriter) () {
	aw.Respond(struct{ Message string `json:"message"` }{Message: "OK"}, http.StatusOK)
}

func (a *Api) NotFoundHandler(ar *ApiRequest, aw *ApiResponseWriter) () {
	aw.Respond(struct{ Message string `json:"message"` }{Message: "Unsupported URL provided."}, http.StatusNotFound)
}

func (a *Api) OperateHandler(operation *Operation) (func(ar *ApiRequest, aw *ApiResponseWriter) ()) {
	return func(ar *ApiRequest, aw *ApiResponseWriter) () {
		type Response OperationResponse

		argumentsQueryParameters := ar.GetQueryParameters()["arguments"]
		if len(argumentsQueryParameters) != 1 {
			err := errors.New(fmt.Sprintf("Could not parse arguments: Invalid number of \"arguments\" query parameters provided. Expected \"1\", but got \"%v\"", len(argumentsQueryParameters)))
			log.Printf(err.Error())
			aw.Respond(&Response{Message: err.Error(), Arguments: make([]float64, 0), Result: 0, AsString: "invalid"}, http.StatusBadRequest)
			return
		}
		argumentsQueryParameter := argumentsQueryParameters[0]

		var argumentsStrings []string
		if argumentsQueryParameter == "" {
			argumentsStrings = make([]string, 0)
		} else {
			argumentsStrings = strings.Split(argumentsQueryParameter, ",")
		}

		if len(argumentsStrings) != operation.ExpectedArgumentLength {
			err := errors.New(fmt.Sprintf("Could not parse arguments: Invalid number of arguments provided in \"arguments\" query parameter. Expected \"%v\", but got \"%v\"", operation.ExpectedArgumentLength, len(argumentsQueryParameters)))
			log.Printf(err.Error())
			aw.Respond(&Response{Message: err.Error(), Arguments: make([]float64, 0), Result: 0, AsString: "invalid"}, http.StatusBadRequest)
			return
		}

		arguments := make([]float64, 0)
		for _, argumentString := range argumentsStrings {
			argument, err := strconv.ParseFloat(argumentString, 64)
			if err != nil {
				err := errors.New(fmt.Sprintf("Could not parse arguments: Could not parse query parameter value \"%v\" into a float64 value.", argumentString))
				log.Printf(err.Error())
				aw.Respond(&Response{Message: err.Error(), Arguments: arguments, Result: 0, AsString: "invalid"}, http.StatusBadRequest)
				return
			}

			arguments = append(arguments, argument)
		}

		result, asString, err := operation.Operate(arguments)
		if err != nil {
			log.Printf(err.Error())
			aw.Respond(&Response{Message: err.Error(), Arguments: arguments, Result: 0, AsString: "invalid"}, http.StatusInternalServerError)
			return
		}

		aw.Respond(&Response{Message: "Success.", Arguments: arguments, Result: result, AsString: asString}, http.StatusOK)
	}
}
