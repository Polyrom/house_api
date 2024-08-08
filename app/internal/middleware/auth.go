package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/Polyrom/houses_api/internal/apierror"
	"github.com/Polyrom/houses_api/pkg/logging"
)

type Token string
type Role string

const UserRole ContextKey = "user_role"
const UserID ContextKey = "user_id"

const (
	Client    Role = "client"
	Moderator Role = "moderator"
)

type isAuthMiddleware struct {
	s Service
	l logging.Logger
}

func (authmw *isAuthMiddleware) DoInMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ContextKeyRequestID).(string)
		token := r.Header.Get("Authorization")
		if token == "" {
			noTokenErr := errors.New("no token")
			authmw.l.Errorf("unauthorized req_id=%s: %v", reqID, noTokenErr)
			apierror.Write(w, noTokenErr, reqID, http.StatusUnauthorized)
			return
		}
		userIDRole, err := authmw.s.GetRoleByToken(r.Context(), Token(token))
		if err != nil {
			authmw.l.Errorf("internal error req_id=%s: %v", reqID, err)
			apierror.Write(w, err, reqID, http.StatusInternalServerError)
			return
		}
		if userIDRole.Role != Client && userIDRole.Role != Moderator {
			noTokenErr := errors.New("not client or moderator")
			authmw.l.Errorf("unauthorized req_id=%s: %v", reqID, noTokenErr)
			apierror.Write(w, noTokenErr, reqID, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserRole, userIDRole.Role)
		ctx = context.WithValue(ctx, UserID, userIDRole.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func NewAuthMiddleware(s Service, l logging.Logger) Middleware {
	return &isAuthMiddleware{s: s, l: l}
}

type isModerMiddleware struct {
	s Service
	l logging.Logger
}

func (modermw *isModerMiddleware) DoInMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(ContextKeyRequestID).(string)
		token := r.Header.Get("Authorization")
		if token == "" {
			noTokenErr := errors.New("no token")
			modermw.l.Errorf("unauthorized req_id=%s: %v", reqID, noTokenErr)
			apierror.Write(w, noTokenErr, reqID, http.StatusUnauthorized)
			return
		}
		userIDRole, err := modermw.s.GetRoleByToken(r.Context(), Token(token))
		if err != nil {
			userNotFoundErr := errors.New("user not found")
			modermw.l.Errorf("internal error req_id=%s: %v", reqID, err)
			apierror.Write(w, userNotFoundErr, reqID, http.StatusInternalServerError)
			return
		}
		if userIDRole.Role != Moderator {
			noTokenErr := errors.New("not a moderator")
			modermw.l.Errorf("unauthorized req_id=%s: %v", reqID, noTokenErr)
			apierror.Write(w, noTokenErr, reqID, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserRole, userIDRole.Role)
		ctx = context.WithValue(ctx, UserID, userIDRole.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func NewIsModerMiddleware(s Service, l logging.Logger) Middleware {
	return &isModerMiddleware{s: s, l: l}
}
