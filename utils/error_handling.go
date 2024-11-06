package utils

type APIError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIErrors struct {
	Errors []APIError `json:"errors"`
}

func NewErrorResponse(field, message string) APIError {
	return APIError{
		Field:   field,
		Message: message,
	}
}
