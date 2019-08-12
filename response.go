package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response is a no-frills standardized response body wrapper with some
// added utility to help manage writing to a http.ResponseWriter.
type Response struct {
	Writer       io.Writer   `json:"-"`
	StatusCode   int         `json:"status_code"`
	StatusText   string      `json:"status_text"`
	ErrorDetails *string     `json:"error_details"`
	Result       interface{} `json:"result"`
}

// New instantiates a new Response struct prepared with a 500 Status.
//
// Internal Server Error is considered a fail-safe to automatically handle
// unexpected panics properly. For everything that is expected, SetResult
// should be called with the proper Status Code before Output is executed.
func New(writer io.Writer) *Response {
	r := new(Response)
	r.Writer = writer
	r.StatusCode = http.StatusInternalServerError
	r.StatusText = http.StatusText(http.StatusInternalServerError)
	return r
}

// WithErrorDetails sets additional human-readable details beyond the
// HTTP Status Code that might help the user consuming this API to
// understand what went wrong, so they can resolve the issue.
//
// WithErrorDetails is entirely optional to call in a Request lifecycle,
// but its usage can lead to a much more consumer friendly API experience.
//
// This tends to be very useful with more ambiguous response codes
// (e.g. 400 Bad Request), but is generally not terribly useful with
// response codes that leave less to question (e.g. 404 Not Found).
//
// For the best code readability, this function should be called by
// being chained onto SetResult:
//
//	resp.SetResult(http.StatusBadRequest, nil).
//		WithErrorDetails("Missing Parameter 'name'")
func (r *Response) WithErrorDetails(errorDetails string) *Response {
	r.ErrorDetails = &errorDetails
	return r
}

// SetResult sets the HTTP Status Code and the Result of this response.
// This should be called once before Output is called, unless you want to
// intentionally throw a 500 response.
func (r *Response) SetResult(httpStatusCode int, result interface{}) *Response {
	r.StatusCode = httpStatusCode
	r.StatusText = http.StatusText(httpStatusCode)
	r.Result = result
	return r
}

// Output writes the Response as a JSON string to the Writer.
//
// If the Writer is a ResponseWriter, we set the Content-Type to
// application/json, and we send an HTTP response header with the
// provided Status Code.
func (r *Response) Output() {
	// Write header, if this is a ResponseWriter
	switch v := r.Writer.(type) {
	case http.ResponseWriter:
		v.Header().Set("Content-Type", "application/json")
		v.WriteHeader(r.StatusCode)
	}

	// Write Body
	b, err := json.Marshal(r)
	if err != nil {
		panic("Unable to json.Marshal our Response")
	}
	fmt.Fprint(r.Writer, string(b))
}
