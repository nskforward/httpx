package response

type APIError struct {
	Status int
	Text   string
}

func (e APIError) Error() string {
	return e.Text
}

func NewAPIError(statusCode int) APIError {
	return APIError{Status: statusCode}
}
