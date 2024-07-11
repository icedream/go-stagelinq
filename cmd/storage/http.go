package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func logMuxHandling(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mux, handling: %s %s", r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	})
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mux, not found: %s %s", r.Method, r.URL.String())
	w.WriteHeader(http.StatusNotFound)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Println("HTTP: Ping")
	w.WriteHeader(http.StatusOK)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestedPath, err := url.PathUnescape(vars["path"])
	if err != nil {
		log.Println("HTTP: Download, bad path:", requestedPath)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if requestedPath != demoTrackURLGRPC {
		log.Println("HTTP: Download, not found:", requestedPath)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println("HTTP: Download, OK:", requestedPath)
	w.Header().Set("Content-type", "application/octet-stream")
	w.Header().Set("Content-length", strconv.Itoa(len(demoTrackBytes)))
	w.WriteHeader(http.StatusOK)
	f := bytes.NewReader(demoTrackBytes)
	io.Copy(w, f)
}

func eaasHTTPHandler() http.Handler {
	r := mux.NewRouter()
	r.Use(logMuxHandling)
	r.NotFoundHandler = http.HandlerFunc(handleNotFound)
	r.UseEncodedPath()
	r.SkipClean(true)
	r.HandleFunc("/download/{path}", handleDownload).Methods(http.MethodGet)
	r.HandleFunc("/ping", handlePing).Methods(http.MethodGet)
	return r
}
