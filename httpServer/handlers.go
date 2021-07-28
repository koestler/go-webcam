package httpServer

import (
	"errors"
	"net/http"
)

func writeJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Model", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func HandleApiNotFound(env *Environment, w http.ResponseWriter, r *http.Request) Error {
	err := errors.New("api method not found")
	return StatusError{404, err}
}

