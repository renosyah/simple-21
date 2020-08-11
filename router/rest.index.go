package router

import (
	"net/http"
)

func (h *RouterHub) HandleIndex(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "max-age=-1")

		// handle query
		// request post
		// etc

		next.ServeHTTP(w, r)
	})
}
