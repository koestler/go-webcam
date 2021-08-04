package httpServer

import "net/http"

type view struct {
	Route       string
	Name        string
	Cameras     []string
	HasHtaccess bool
}

func HandleViewsIndex(env *Environment, w http.ResponseWriter, r *http.Request) Error {
	views := make([]view, 0)

	for _, v := range env.Views {
		views = append(views, view{
			Route:       v.Route(),
			Name:        v.Name(),
			Cameras:     v.Cameras(),
			HasHtaccess: false,
		})
	}

	return writeJsonResponse(w, views)
}
