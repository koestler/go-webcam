package httpServer

import (
	"expvar"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/apache-logformat"
	"io"
	"log"
	"net/http"
)

type HttpRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerHandleFunc
}

var staticHttpRoutes = []HttpRoute{
	{
		"views",
		"GET",
		"/api/v0/views",
		HandleViewsIndex,
	}, {
		"expvar",
		"GET",
		"/debug/vars",
		func(env *Environment, w http.ResponseWriter, r *http.Request) Error {
			expvar.Handler().ServeHTTP(w, r)
			return nil
		},
	},
}

func newRouter(logger io.Writer, env *Environment) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	httpRoutes := append(staticHttpRoutes, getDynamicHttpRoutes(env)...)

	// setup normal http routes
	for _, route := range httpRoutes {
		var handler http.Handler
		handler = Handler{Env: env, Handle: route.HandlerFunc}
		if logger != nil {
			handler = apachelog.CombinedLog.Wrap(handler, logger)
		}

		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

		log.Printf("httpServer: setup route for: %v %v : %v", route.Method, route.Pattern, route.Name)
	}

	return router
}

func getDynamicHttpRoutes(env *Environment) []HttpRoute {
	routes := make([]HttpRoute, 0)

	for _, v := range env.Views {
		view := v
		routes = append(routes, HttpRoute{
			"index page",
			"GET",
			"/" + view.Name(),
			func(env *Environment, w http.ResponseWriter, r *http.Request) Error {
				return handleViewIndex(view, w, r)
			},
		})
		for _, c := range view.Cameras() {
			camera := c

			cameraClient := env.CameraClientPoolInstance.GetClient(camera)
			if cameraClient == nil {
				continue
			}

			routes = append(routes, HttpRoute{
				"fetch image",
				"GET",
				"/api/v0/images/" + view.Name() + "/" + camera + ".jpg",
				func(env *Environment, w http.ResponseWriter, r *http.Request) Error {
					return handleCameraImage(cameraClient, view, w, r)
				},
			})
		}

	}

	return routes
}
