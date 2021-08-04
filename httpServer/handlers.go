package httpServer

import (
	"encoding/json"
	"net/http"
)

func writeJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Model", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func writeJsonResponse(w http.ResponseWriter, data interface{}) Error {
	writeJsonHeaders(w)
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return StatusError{500, err}
	}
	_, err = w.Write(b)
	if err != nil {
		return StatusError{500, err}
	}
	return nil
}
