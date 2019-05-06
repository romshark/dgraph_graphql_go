package http

// GraphResponseError represents a response error object
type GraphResponseError struct {
	Code    string `json:"c,omitempty"`
	Message string `json:"m,omitempty"`
}

// GraphResponse represents a response object
type GraphResponse struct {
	Data  interface{}         `json:"d,omitempty"`
	Error *GraphResponseError `json:"e,omitempty"`
}
