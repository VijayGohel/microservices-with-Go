package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON reads and decodes JSON data from the request body into the provided data interface.
// It also validates that the request body contains a single JSON value and does not exceed the maximum allowed size.
//
// Parameters:
// - w: http.ResponseWriter to write any error response.
// - r: *http.Request containing the request body to read from.
// - data: any interface where the decoded JSON data will be stored.
//
// Returns:
// - error: nil if the JSON data is successfully read and decoded, or an error if any issues occur during the process.
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 //1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have single json value")
	}

	return nil
}

// writeJSON writes JSON data to the HTTP response with the specified status code and optional headers.
// It marshals the provided data into JSON format and writes it to the response body.
// If any errors occur during the process, such as marshaling or writing to the response, an error is returned.
//
// Parameters:
// - w: http.ResponseWriter to write the response.
// - status: int representing the HTTP status code to set in the response.
// - data: any interface containing the data to be marshalled into JSON format.
// - headers: ...http.Header representing optional headers to set in the response.
//
// Returns:
// - error: nil if the JSON data is successfully written to the response, or an error if any issues occur during the process.
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON writes an error response in JSON format to the HTTP response with the specified status code.
// It constructs a jsonResponse struct with the error flag set to true and the provided error message.
// Then, it calls the writeJSON function to write the jsonResponse to the HTTP response.
//
// Parameters:
// - w: http.ResponseWriter to write the response.
// - err: error containing the error message to be included in the response.
// - status: int representing the HTTP status code to set in the response.
//
// Returns:
// - error: nil if the JSON data is successfully written to the response, or an error if any issues occur during the process.
func (app *Config) errorJSON(w http.ResponseWriter, err error, status int) error {
	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return app.writeJSON(w, status, payload)
}
