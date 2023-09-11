package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) handleSlackDetails(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetSlackDetails(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api", makeHTTPHandleFunc(s.handleSlackDetails))

	log.Println("Rest API running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}

func (s *APIServer) handleGetSlackDetails(w http.ResponseWriter, r *http.Request) error {
	return nil
}
