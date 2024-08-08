package house

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Polyrom/houses_api/internal/apierror"
	"github.com/Polyrom/houses_api/internal/handlers"
	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const (
	createURL    = "/house/create"
	subscribeURL = "/house/{id}/subscribe"
)

type handler struct {
	aumw  middleware.Middleware
	modmw middleware.Middleware
	s     *Service
	l     logging.Logger
}

func NewHandler(aumw middleware.Middleware, modmw middleware.Middleware, s *Service, l logging.Logger) handlers.Handler {
	return &handler{aumw: aumw, modmw: modmw, s: s, l: l}
}

func (h *handler) Register(r *mux.Router) {
	r.Handle(createURL, h.modmw.DoInMiddle(http.HandlerFunc(h.Create))).Methods(http.MethodPost)
	r.Handle(subscribeURL, h.aumw.DoInMiddle(http.HandlerFunc(h.Subscribe))).Methods(http.MethodPost)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	// FIXME: time is all zeroes
	reqID := r.Context().Value(middleware.ContextKeyRequestID).(string)
	var hdto CreateHouseDTO
	err := json.NewDecoder(r.Body).Decode(&hdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(hdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	newHouse, err := h.s.Create(r.Context(), hdto)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newHouse)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
}

func (h *handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	// FIXME: implement
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "house created")
}
