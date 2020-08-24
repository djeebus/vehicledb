package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var cache = make(map[string]*gojsonschema.Schema)

func validateSchemaBuildModel(request *http.Request, schemaText string, model interface{}) error {
	var err error

	schema, ok := cache[schemaText]
	if !ok {
		loader := gojsonschema.NewStringLoader(schemaText)
		sl := gojsonschema.NewSchemaLoader()
		schema, err = sl.Compile(loader)
		if err != nil {
			return err
		}

		cache[schemaText] = schema
	}

	buf, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}

	document := gojsonschema.NewBytesLoader(buf)
	result, err := schema.Validate(document)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return &ValidationError{result}
	}

	buffer := bytes.NewBuffer(buf)
	decoder := json.NewDecoder(buffer)
	err = decoder.Decode(&model)
	if err != nil {
		return err
	}

	return nil
}

type ValidationError struct {
	Result *gojsonschema.Result
}

func (error ValidationError) Error() string {
	errors := make([]string, len(error.Result.Errors()))

	for _, err := range error.Result.Errors() {
		errors = append(errors, err.String())
	}
	return strings.Join(errors, "\n")
}

func renderJson(writer io.Writer, data interface{}) {
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		fmt.Println("failed to write json")
	}
}

func renderError(writer http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *ValidationError:
		writer.WriteHeader(400)

		errors := make([]interface{}, 0)

		for _, result := range e.Result.Errors() {
			switch r := result.(type) {
			case *gojsonschema.RequiredError:
				errors = append(errors, map[string]interface{}{"code": "missing_field", "field": r.Details()["property"]})
			default:
				errors = append(errors, map[string]interface{}{"code": r.Type()})
			}

		}
		renderJson(writer, map[string]interface{}{
			"code":   "invalid_request",
			"errors": errors,
		})

	default:
		fmt.Println(err)
		writer.WriteHeader(500) // todo: use 400 for jsonschema errors
		renderJson(writer, map[string]interface{}{"code": "server_error"})
	}
}

func keys(mapped map[string]http.HandlerFunc) []string {
	returnable := make([]string, 0, len(mapped))
	for key := range mapped {
		returnable = append(returnable, key)
	}
	return returnable
}

func AddMappedMethods(route *mux.Route, methodMap map[string]http.HandlerFunc) {
	route.Methods(keys(methodMap)...).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methodMap[r.Method](w, r)
	})
}
