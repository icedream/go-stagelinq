package main

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func newHTTPServiceHandler() http.Handler {
	r := mux.NewRouter()
	r.Get("/download/{path}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestedPath := vars["path"]
		if requestedPath != demoTrackURL {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-length", strconv.Itoa(len(demoTrackBytes)))
		w.WriteHeader(http.StatusOK)
		f := bytes.NewReader(demoTrackBytes)
		io.Copy(w, f)
	})
	r.Get("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return r
}
