package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Polyrom/houses_api/internal/apierror"
	"github.com/Polyrom/houses_api/internal/handlers"
	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const (
	loginURL          = "/login"
	registerURL       = "/register"
	dummyLoginURL     = "/dummyLogin"
	dummyEmailSuffix  = "@foo.huh"
	dummyUserPassword = "dummyPass"
)

type handler struct {
	s *Service
	l logging.Logger
}

func NewHandler(s *Service, l logging.Logger) handlers.Handler {
	return &handler{s, l}
}

func (h *handler) Register(r *mux.Router) {
	r.HandleFunc(loginURL, h.UserLogin).Methods(http.MethodPost)
	r.HandleFunc(registerURL, h.UserRegister).Methods(http.MethodPost)
	r.HandleFunc(dummyLoginURL, h.UserDummyLogin).Methods(http.MethodGet)
}

func (h *handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.ContextKeyRequestID).(string)
	var userdto UserRegisterDTO
	err := json.NewDecoder(r.Body).Decode(&userdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(userdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	user := User{
		Email:    userdto.Email,
		Password: userdto.Password,
		Role:     userdto.Role,
	}
	userid, err := h.s.Register(r.Context(), user)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	useridResp := UserIDDTO{UserID: userid}
	err = json.NewEncoder(w).Encode(useridResp)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
}

func (h *handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.ContextKeyRequestID).(string)
	var uldto UserLoginDTO
	err := json.NewDecoder(r.Body).Decode(&uldto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	storedUser, err := h.s.GetByID(r.Context(), uldto.UserID)
	if err != nil {
		notFoundErr := errors.New("user not found")
		h.l.Errorf("user not found req_id=%s: %v", reqID, err)
		apierror.Write(w, notFoundErr, reqID, http.StatusNotFound)
		return
	}
	err = storedUser.VerifyPassword(uldto.Password)
	if err != nil {
		passwErr := errors.New("wrong password")
		h.l.Errorf("wrong password req_id=%s: %v", reqID, err)
		apierror.Write(w, passwErr, reqID, http.StatusUnauthorized)
		return
	}
	token := storedUser.GenerateToken()
	err = h.s.AddToken(r.Context(), storedUser.ID, token)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	useridResp := TokenDTO{Token: token}
	err = json.NewEncoder(w).Encode(useridResp)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
}

func (h *handler) UserDummyLogin(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.ContextKeyRequestID).(string)
	params := r.URL.Query()
	userType := middleware.Role(params.Get("user_type"))
	if userType == middleware.Role("") {
		noUserTypeErr := errors.New("user type not specified")
		h.l.Errorf("user not found req_id=%s: %v", reqID, noUserTypeErr)
		apierror.Write(w, noUserTypeErr, reqID, http.StatusNotFound)
		return
	}
	if userType != middleware.Moderator && userType != middleware.Client {
		unknownUserTypeErr := errors.New("unknown user type (client or moderator)")
		h.l.Errorf("user not found req_id=%s: %v", reqID, unknownUserTypeErr)
		apierror.Write(w, unknownUserTypeErr, reqID, http.StatusNotFound)
		return
	}
	dummyEmailPrefix := h.s.GenerateRandomEmailPrefix(r.Context(), 20)
	dummyUser := User{
		Email:    dummyEmailPrefix + dummyEmailSuffix,
		Password: dummyUserPassword,
		Role:     string(userType),
	}
	dummyUserID, err := h.s.Register(r.Context(), dummyUser)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	token := dummyUser.GenerateToken()
	err = h.s.AddToken(r.Context(), dummyUserID, token)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	useridResp := TokenDTO{Token: token}
	err = json.NewEncoder(w).Encode(useridResp)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
}
