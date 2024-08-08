package flat

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Polyrom/houses_api/internal/apierror"
	"github.com/Polyrom/houses_api/internal/handlers"
	"github.com/Polyrom/houses_api/internal/middleware"
	"github.com/Polyrom/houses_api/internal/modstatus"
	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const (
	createURL   = "/flat/create"
	updateURL   = "/flat/update"
	findByIDURL = "/house/{id}"
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
	r.Handle(createURL, h.aumw.DoInMiddle(http.HandlerFunc(h.Create))).Methods(http.MethodPost)
	r.Handle(updateURL, h.modmw.DoInMiddle(http.HandlerFunc(h.Update))).Methods(http.MethodPost)
	r.Handle(findByIDURL, h.aumw.DoInMiddle(http.HandlerFunc(h.FindByID))).Methods(http.MethodGet)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.ContextKeyRequestID).(string)
	var fdto CreateFlatDTO
	err := json.NewDecoder(r.Body).Decode(&fdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(fdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	newFlat, err := h.s.Create(r.Context(), fdto)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(newFlat)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.ContextKeyRequestID).(string)
	var ufsdto UpdateFlatStatusDTO
	err := json.NewDecoder(r.Body).Decode(&ufsdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	validStatuses := []string{modstatus.Approved.String(), modstatus.Declined.String(), modstatus.OnModeration.String()}
	validate := validator.New()
	validate.RegisterValidation("oneof_modstat", func(fl validator.FieldLevel) bool {
		for _, allowed := range validStatuses {
			if fl.Field().String() == allowed {
				return true
			}
		}
		return false
	})
	err = validate.Struct(ufsdto)
	if err != nil {
		h.l.Errorf("bad request req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusBadRequest)
		return
	}
	updatedFlat, err := h.s.Update(r.Context(), ufsdto)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedFlat)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
}

func (h *handler) FindByID(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value(middleware.ContextKeyRequestID).(string)
	var fid FlatID
	vars := mux.Vars(r)
	fidParam, ok := vars["id"]
	if !ok {
		noIDErr := errors.New("flat id not found")
		h.l.Errorf("bad request req_id=%s: %v", reqID, noIDErr)
		apierror.Write(w, noIDErr, reqID, http.StatusBadRequest)
		return
	}
	fidParamConv, err := strconv.Atoi(fidParam)
	if err != nil {
		invalidIDErr := errors.New("invalid flat id")
		h.l.Errorf("bad request req_id=%s: %v", reqID, invalidIDErr)
		apierror.Write(w, invalidIDErr, reqID, http.StatusBadRequest)
		return
	}
	fid = FlatID(fidParamConv)
	flatsFound, err := h.s.GetByHouseID(r.Context(), fid)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(flatsFound)
	if err != nil {
		h.l.Errorf("internal error req_id=%s: %v", reqID, err)
		apierror.Write(w, err, reqID, http.StatusInternalServerError)
		return
	}
}
