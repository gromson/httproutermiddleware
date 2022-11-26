package httprouter_middleware

import "github.com/julienschmidt/httprouter"

type Routes []Route

// Route struct for representing a route.
type Route struct {
	// Path
	Path string
	// Handle function
	Handler httprouter.Handle
	// Method verb
	Method string
	// Set of middlewares to be applied to the routes
	Pipeline Pipeline
}
