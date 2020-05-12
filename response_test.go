package response

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

// StackWriter is a Writer that allows us to Peek to see what was
// last response written to it. StackWriter is also a http.ResponseWriter.
type StackWriter struct {
	HeaderInt int
	Stack     []string
}

func (sw *StackWriter) Write(p []byte) (n int, err error) {
	sw.Stack = append(sw.Stack, string(p))
	return 0, nil
}

func (sw *StackWriter) WriteHeader(header int) {
	sw.HeaderInt = header
}

func (sw *StackWriter) Header() http.Header {
	return make(http.Header, 0)
}

func (sw *StackWriter) Peek() *Response {
	if len(sw.Stack) > 0 {
		popped := &Response{}
		err := json.Unmarshal([]byte(
			sw.Stack[len(sw.Stack)-1],
		), popped)
		if err != nil {
			panic(err)
		}
		return popped
	}
	return nil
}

// TestResponseNotSet tests our response.New default functionality
func TestResponseDefault(t *testing.T) {
	// Prepare StackWriter
	sw := new(StackWriter)

	// Prepare Response
	resp := New(sw)
	resp.Output()

	// Validate
	assert.Equal(t, Response{
		StatusCode:   http.StatusInternalServerError,
		StatusText:   "Internal Server Error",
		ErrorDetails: nil,
		Result:       nil,
	}, *(sw.Peek()))
}

// TestResponseWithErrorDetails tests our WithErrorDetails functionality
func TestResponseWithErrorDetails(t *testing.T) {
	// Prepare StackWriter
	sw := new(StackWriter)

	// Prepare Response
	resp := New(sw)
	resp.SetResult(http.StatusInternalServerError, nil).
		WithErrorDetails("Missing Auth")
	resp.Output()

	// Validate
	missingAuth := "Missing Auth"
	assert.Equal(t, Response{
		StatusCode:   http.StatusInternalServerError,
		StatusText:   "Internal Server Error",
		ErrorDetails: &[]string{missingAuth},
		Result:       nil,
	}, *(sw.Peek()))
}

// TestResponseSuccess tests a successful response
func TestResponseSuccess(t *testing.T) {
	// Prepare StackWriter
	sw := new(StackWriter)

	// Prepare Response
	resp := New(sw)
	resp.SetResult(http.StatusOK,
		struct {
			ValueOne string `json:"value_one"`
			ValueTwo string `json:"value_two"`
		}{
			ValueOne: "foo",
			ValueTwo: "bar",
		},
	)
	resp.Output()

	// Validate
	assert.Equal(t, Response{
		StatusCode:   http.StatusOK,
		StatusText:   "OK",
		ErrorDetails: nil,
		Result: map[string]interface{}{
			"value_one": "foo",
			"value_two": "bar",
		},
	}, *(sw.Peek()))
}

// TestJsonRenderFailure tests the scenario when the input fails to be
// rendered to JSON, and ensures that a panic is called.
func TestJsonRenderFailure(t *testing.T) {
	defer func() {
		recover()
	}()
	sw := new(StackWriter)
	resp := New(sw)
	resp.SetResult(http.StatusOK, func() {})
	resp.Output()
	t.Error("JsonRenderer should fail and cause a panic with content that can not be serialized to JSON")
}

// TestParseSuccess tests that we can parse a properly formatted
// io.ReadCloser into a Response struct
func TestParseSuccess(t *testing.T) {
	body := ioutil.NopCloser(
		bytes.NewReader([]byte(
			`{"status_code":200,"status_text":"OK","error_details":null,"result":"pong"}`,
		)),
	)
	result := Parse(body)
	assert.Equal(t, Response{
		StatusCode:   200,
		StatusText:   "OK",
		ErrorDetails: nil,
		Result:       "pong",
	}, result)
}

// TestParseFailure tests that a panic is properly thrown
// when we try to Parse a io.ReadCloser is improperly formatted
func TestParseFailure(t *testing.T) {
	defer func() {
		recover()
	}()
	body := ioutil.NopCloser(
		bytes.NewReader([]byte(
			`{"status_code":"hello world"}`,
		)),
	)
	Parse(body)
	t.Error("Parse should fail and cause a panic with content that can not be converted to a Response")
}
