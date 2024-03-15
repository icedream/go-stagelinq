package main

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
)

func newHTTPServiceHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/download/{path}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.PathValue("path")
		if requestedPath != demoTrackURL {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-length", strconv.Itoa(len(demoTrackBytes)))
		w.WriteHeader(http.StatusOK)
		f := bytes.NewReader(demoTrackBytes)
		io.Copy(w, f)
	}))
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	return mux
}
