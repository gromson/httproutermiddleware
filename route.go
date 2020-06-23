package httprouter_middleware

import "github.com/julienschmidt/httprouter"

// Route struct for representing a route
type Route struct {
	// Path
	Path string
	// Handle function
	Handler httprouter.Handle
	// Method verb
	Method string
	// Set of middlewares to be applied to the routes
	Pipeline []Middleware
}

func (r *Route) wrap() httprouter.Handle {
	return wrap(r.Handler, r.Pipeline)
}
