package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rivermq/rivermq/util"
)

// NewRiverMQRouter does something
func NewRiverMQRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = util.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
		if route.Queries != nil {
			for key := range route.Queries {
				router.Queries(key, route.Queries[key])
			}
		}
	}
	return router
}
