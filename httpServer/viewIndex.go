package httpServer

import (
	"fmt"
	"github.com/koestler/go-webcam/config"
	"net/http"
)

func handleViewIndex(view *config.ViewConfig, w http.ResponseWriter, r *http.Request) Error {
	fmt.Fprintf(w, "<h1>%s</h1>", view.Name())

	fmt.Fprintln(w, "<ul>")

	for _, camera := range view.Cameras() {
		fmt.Fprintf(w, "<li><img src=\"/api/v0/images/%s/%s.jpg?width=400\" /></li>", view.Name(), camera)
	}

	fmt.Fprintln(w, "</ul>")

	return nil
}
