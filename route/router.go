package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rivermq/rivermq/handler"
	"github.com/rivermq/rivermq/util"
)

// NewRiverMQRouter creates a mux.Router configured with the various
// handlers of this application
func NewRiverMQRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	registerHandler(router,
		"CreateSubscriptionHandler",
		"POST",
		"/subscriptions",
		handler.CreateSubscriptionHandler,
		nil)

	registerHandler(router,
		"FindSubscriptionByIDHandler",
		"GET",
		"/subscriptions/{subID}",
		handler.FindSubscriptionByIDHandler,
		nil)

	registerHandler(router,
		"FindAllSubscriptionsHandler",
		"GET",
		"/subscriptions",
		handler.FindAllSubscriptionsHandler,
		map[string]string{
			"type": "{type:^\\w+$}",
		})

	registerHandler(router,
		"DeleteSubscriptionByIDHandler",
		"DELETE",
		"/subscriptions/{subID}",
		handler.DeleteSubscriptionByIDHandler,
		nil)

	registerHandler(router,
		"CreateMessageHander",
		"POST",
		"/messages",
		handler.CreateMessageHander,
		nil)

	return router
}

func registerHandler(router *mux.Router,
	name string,
	method string,
	path string,
	handlerFunc http.HandlerFunc,
	queryParams map[string]string) {

	handler := util.Logger(handlerFunc, name)
	router.
		Methods(method).
		Path(path).
		Name(name).
		Handler(handler)

	if queryParams != nil {
		for key := range queryParams {
			router.Queries(key, queryParams[key])
		}
	}
}
