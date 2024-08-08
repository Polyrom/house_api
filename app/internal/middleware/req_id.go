package middleware

import (
	"context"
	"net/http"

	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/google/uuid"
)

type ContextKey string

const ContextKeyRequestID ContextKey = "requestID"

type reqIDMiddleware struct {
	l logging.Logger
}

func (ridmw *reqIDMiddleware) DoInMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()
		ctx = context.WithValue(ctx, ContextKeyRequestID, id.String())
		r = r.WithContext(ctx)
		ridmw.l.Infof("request %s %s req_id=%s", r.Method, r.URL, ctx.Value(ContextKeyRequestID))
		next.ServeHTTP(w, r)
		ridmw.l.Infof("request %s %s handled req_id=%s", r.Method, r.URL, ctx.Value(ContextKeyRequestID))
	})
}

func NewReqIDMiddleware(l logging.Logger) Middleware {
	return &reqIDMiddleware{l: l}
}
