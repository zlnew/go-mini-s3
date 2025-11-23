// Package middleware
package middleware

import (
	"net/http"
)

type Middleware = func(http.Handler) http.Handler

func Chain(mw ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(mw) - 1; i >= 0; i-- {
			final = mw[i](final)
		}
		return final
	}
}

func UsePublic() Middleware {
	handler := Chain(Logger, Recoverer)
	return handler
}

func UseProtected() Middleware {
	handler := Chain(Logger, Recoverer, APIKey)
	return handler
}
