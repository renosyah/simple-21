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
		Status       int `json:"status"`
		BaseResponse `json:"errors"`
		Data         interface{} `json:"result"`
	}
)

func HttpResponse(w http.ResponseWriter, r *http.Request, data interface{}, status int) {

	resp := Response{
		Status: status,
		Data:   data,
		BaseResponse: BaseResponse{
			Errors: []string{},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func HttpResponseException(w http.ResponseWriter, r *http.Request, status int) {

	resp := Response{
		Status: status,
		Data:   nil,
		BaseResponse: BaseResponse{
			Errors: []string{},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
