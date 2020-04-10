package base_api

import (
	"errors"
)

// This type is used to signify to handlers that the error message is okay and actually intended to be show to the end user, while still allowing for all functions to return the generic "error" type.
// This is done by having this struct implement the error interface (hence it is of the generic type "error"), but in the handler we check if the type assertion to this more specific type is valid.
// 		err := doThing()
// 		if apiErr, isApiErr := models.IsApiError(err); isApiErr { /* apiErr contains an error message which we can display to the end user */ }
// Note that ApiError is simply a wrapper around the generic error type.
type ApiError struct {
	error
}

// Required function for ApiError to implement the "error" interface.
func (a ApiError) Error() (message string) {
	return a.error.Error()
}

func NewApiError(message string) (apiError *ApiError) {
	return &ApiError{error: errors.New(message)}
}

func IsApiError(err error) (apiErr *ApiError, isApiError bool) {
	apiErr, isApiError = err.(*ApiError)
	return apiErr, isApiError
}
