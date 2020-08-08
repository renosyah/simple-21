package api

import (
	"encoding/json"
	"net/http"
)

type (
	BaseResponse struct {
		Errors []string `json:"errors,omitempty"`
	}
	Response struct {
		Status       string `json:"status"`
		BaseResponse `json:"errors"`
		Data         interface{} `json:"result"`
	}
)

func HttpResponse(w http.ResponseWriter, r *http.Request, data interface{}, status int) {

	resp := Response{
		Status: http.StatusText(status),
		Data:   data,
		BaseResponse: BaseResponse{
			Errors: []string{},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
