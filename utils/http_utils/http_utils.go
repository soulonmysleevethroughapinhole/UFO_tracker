package http_utils

import (
	"encoding/json"
	"net/http"

	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

func RespondJson(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func RespondError(w http.ResponseWriter, err rest_errors.RestErr) {
	RespondJson(w, err.Status(), err)
}
