package middleware

import "github.com/wit-switch/assessment-tax/config"

type Dependencies struct {
	Auth *config.AuthConfig
}

type Middleware struct {
	auth *config.AuthConfig
}

func NewMiddleware(deps Dependencies) *Middleware {
	return &Middleware{
		auth: deps.Auth,
	}
}
