package httpServer

import "net/http"

type view struct {
	Name        string
	Title       string
	Cameras     []string
	HasHtaccess bool
}

func HandleViewsIndex(env *Environment, w http.ResponseWriter, r *http.Request) Error {
	views := make([]view, 0)

	for _, v := range env.Views {
		views = append(views, view{
			Name:        v.Name(),
			Title:       v.Title(),
			Cameras:     v.Cameras(),
			HasHtaccess: false,
		})
	}

	return writeJsonResponse(w, views)
}
