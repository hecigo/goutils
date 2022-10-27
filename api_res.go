package goutils

// The default payload struct for API response.
type APIRes struct {
	// The same HTTP status code as the response.
	Status int `json:"status"`

	// A specific code to identify the error.
	ErrorCode string `json:"error_code"`

	// Error message.
	Message string `json:"message"`

	// The data payload.
	Data interface{} `json:"data"`
}

// Makes it compatible with the `error` interface.
func (r *APIRes) Error() string {
	if r.Status >= 400 {
		return r.Message
	}
	return ""
}
