package route

import (
	"net/http"

	"github.com/rivermq/rivermq/handler"
)

// Route defines the Route structure
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Queries     map[string]string
}

// Routes is a slice of Routes
type Routes []Route

var routes = Routes{
	Route{
		"CreateSubscriptionHandler",
		"POST",
		"/subscriptions",
		handler.CreateSubscriptionHandler,
		nil,
	},
	Route{
		"FindSubscriptionByIDHandler",
		"GET",
		"/subscriptions/{subID}",
		handler.FindSubscriptionByIDHandler,
		nil,
	},
	Route{
		"FindAllSubscriptionsHandler",
		"GET",
		"/subscriptions",
		handler.FindAllSubscriptionsHandler,
		nil,
	},
	Route{
		"",
		"GET",
		"/subscriptions",
		handler.FindAllSubscriptionsHandler,
		map[string]string{
			"type": "{type:^\\w+$}",
		},
	},
	Route{
		"DeleteSubscriptionByIDHandler",
		"DELETE",
		"/subscriptions/{subID}",
		handler.DeleteSubscriptionByIDHandler,
		nil,
	},
}
