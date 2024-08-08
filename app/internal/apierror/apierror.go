package apierror

import (
	"encoding/json"
	"net/http"
)

type body struct {
	Message string `json:"message"`
	ReqID   string `json:"req_id"`
	ErrCode int    `json:"err_code"`
}

func Write(w http.ResponseWriter, APIErr error, reqID string, errCode int) {
	w.WriteHeader(errCode)
	w.Header().Set("Retry-After", "5")
	errBody := body{
		Message: APIErr.Error(),
		ReqID:   reqID,
		ErrCode: errCode,
	}
	_ = json.NewEncoder(w).Encode(errBody)
}
