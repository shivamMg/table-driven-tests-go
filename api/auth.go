package api

type Authenticator interface {
	IsAuthenticated(token string) bool
}

type AuthClient struct{}

func (cli *AuthClient) IsAuthenticated(token string) bool {
	return true
}
