package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

/*
ParseJSON decodes a JSON request body into the provided payload.

@param r - *http.Request: The HTTP request containing the JSON body.
@param payload - any: A pointer to the struct where the JSON data will be stored.

@return error - Returns an error if the request body is missing or cannot be parsed.
*/
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

/*
WriteJSON encodes a given value into JSON and writes it to the HTTP response.

@param w - http.ResponseWriter: The response writer to send the JSON response.
@param status - int: The HTTP status code for the response.
@param v - any: The value to encode and send as a JSON response.

@return error - Returns an error if encoding fails.
*/
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

/*
WriteError sends a JSON error response with a given HTTP status code.

@param w - http.ResponseWriter: The response writer to send the error response.
@param status - int: The HTTP status code for the error response.
@param err - error: The error message to include in the response.
*/
func WriteError(w http.ResponseWriter, status int, err error) {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}
