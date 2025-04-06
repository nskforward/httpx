package httpx

type APIError struct {
	Code   int
	Mesage string
}

func (apiError *APIError) Error() string {
	return apiError.Mesage
}

func Throw(code int, msg string) *APIError {
	return &APIError{
		Code:   code,
		Mesage: msg,
	}
}
