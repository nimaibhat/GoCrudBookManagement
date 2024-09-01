package response

import (
	"encoding/json"
	"net/http"
)

// Define a response structure
type Response struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
}

// Helper function to create a new response
func NewResponse(data interface{}, status int) *Response {
	return &Response{
		Status: status,
		Result: data,
	}
}

// Send the response as JSON
func (resp *Response) SendResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	json.NewEncoder(w).Encode(resp)
}

// Convert response to JSON string
func (resp *Response) String() string {
	data, _ := json.Marshal(resp)
	return string(data)
}
