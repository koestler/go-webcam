package httpServer

import (
	"fmt"
	"net/http"
)

func handleCameraImage(camera string, w http.ResponseWriter, r *http.Request) Error {
	fmt.Fprintf(w, "<h1>%s</h1>", camera)

	return nil
}
