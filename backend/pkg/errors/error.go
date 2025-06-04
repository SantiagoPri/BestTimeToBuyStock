package errors

// Error represents a standard error response
// @Description Standard error response format
type Error struct {
	// @Description The error code identifying the type of error
	// @Example validation_error
	Code string `json:"code"`

	// @Description A human-readable error message
	// @Example Invalid input parameters
	Message string `json:"message"`

	// @Description Optional details about the error
	// @Example Stack trace or additional context
	Details string `json:"details,omitempty"`
}
