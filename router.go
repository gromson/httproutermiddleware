package httprouter_middleware

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Middleware wraps the next handler
type Middleware func(next httprouter.Handle) httprouter.Handle

// httprouter.Router wrapper
type Router struct {
	// Groups
	Groups []Group
	// Routes
	Routes []Route
	*httprouter.Router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.apply()
	r.Router.ServeHTTP(w, req)
}

func (r *Router) apply() {
	r.applyRoutes()
	r.applyGroups()
}

func (r *Router) applyRoutes() {
	for _, route := range r.Routes {
		h := route.wrap()
		r.Handle(route.Method, route.Path, h)
	}
}

func (r *Router) applyGroups() {
	for _, group := range r.Groups {
		group.apply(r)
	}
}
