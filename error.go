package witty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"error"` // error message
	Code     string         `json:"code"`  // error code
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		e.Response.Request.Method, e.Response.Request.URL,
		e.Response.StatusCode, e.Message, e.Code)
}

// CheckResponse checks the API response for error, and returns it if
// present.
func CheckResponse(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	errorResponse := &ErrorResponse{Response: resp}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}
