package middleware

import "net/http"

type Middleware interface {
	DoInMiddle(next http.Handler) http.Handler
}
