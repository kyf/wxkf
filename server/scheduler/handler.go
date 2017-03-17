package main

import (
	"net/http"
)

type Handler struct{}

func (this *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	process := getBestProcessor()

	if process == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(process.IP))
}
