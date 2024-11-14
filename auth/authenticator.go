package auth

type Authenticator interface {
	Authenticate(string) (any, error)
}
