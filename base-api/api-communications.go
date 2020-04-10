package base_api

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"errors"
	"encoding/json"
	"strconv"
	"net/url"
	"time"
)

type ApiRequest struct {
	*http.Request
}

func NewApiRequest(r *http.Request) (ar *ApiRequest) {
	ar = new(ApiRequest)
	ar.Request = r
	log.Printf("Call received: \"%v %v\"\n", ar.Method, ar.URL.Path)
	return ar
}

func (ar *ApiRequest) GetRouteVariables() (routeVariables map[string]string) {
	return mux.Vars(ar.Request)
}

func (ar *ApiRequest) GetQueryParameters() (queryParameters url.Values) {
	return ar.URL.Query()
}

func (ar *ApiRequest) GetRequestBody(destination interface{}) (err error) {
	rawRequestBody, err := ioutil.ReadAll(ar.Body)
	if err != nil {
		log.Printf("Request Error: \"%v\"\n", err.Error())
		return errors.New("Could not read request body.\n")
	}
	err = json.Unmarshal(rawRequestBody, &destination)
	if err != nil {
		log.Println(err)
		return errors.New("Could not parse request body.\n")
	}

	return nil
}

func (ar *ApiRequest) GetHeader(headerName string) (headerValue string) {
	return ar.Header.Get(headerName)
}

type ApiResponseWriter struct {
	http.ResponseWriter
	*http.Request
}

func NewApiResponseWriter(w http.ResponseWriter, r *http.Request) (aw *ApiResponseWriter) {
	aw = new(ApiResponseWriter)
	aw.ResponseWriter = w
	aw.Request = r
	return aw
}

func (aw *ApiResponseWriter) Respond(response interface{}, responseStatus int) () {
	aw.ResponseWriter.Header().Set("Content-Type", "application/json")

	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("JSON Marshal Error: \"%v\"\n", err.Error())
		aw.WriteHeader(http.StatusInternalServerError)
		aw.ResponseWriter.Write([]byte("{\n\t\"error\": \"Could not process response body.\"\n}"))
		log.Printf("Response sent: %v: \"%v %v\"\n", strconv.Itoa(http.StatusInternalServerError), aw.Request.Method, aw.Request.URL.Path)
		return
	}

	aw.WriteHeader(responseStatus)
	aw.ResponseWriter.Write([]byte(responseBody))
	log.Printf("Response sent: %v: \"%v %v\"\n", strconv.Itoa(responseStatus), aw.Request.Method, aw.Request.URL.Path)
}

func NewApiCommunication(r *http.Request, w http.ResponseWriter) (ar *ApiRequest, aw *ApiResponseWriter) {
	return NewApiRequest(r), NewApiResponseWriter(w, r)
}

func WrapHandler(apiHandler func(ar *ApiRequest, aw *ApiResponseWriter) (), metrics *PrometheusMetrics) (handler func(w http.ResponseWriter, r *http.Request) ()) {
	return func(w http.ResponseWriter, r *http.Request) () {
		start := time.Now()
		defer func() () {
			end := time.Since(start)
			metrics.Summary.Observe(float64(end))
		}()
		metrics.Counter.Inc()

		ar, aw := NewApiCommunication(r, w)
		apiHandler(ar, aw)
	}
}
