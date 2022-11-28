package httproutermiddleware

import (
	"github.com/julienschmidt/httprouter"
)

// Middleware wraps the next handler.
type Middleware func(next httprouter.Handle) httprouter.Handle

// Pipeline the pipeline of middlawares.
type Pipeline []Middleware

type Config struct {
	BasePath string
	Groups   Groups
	Routes   Routes
	Pipeline Pipeline
}

// Router wrapper for httprouter.Router.
type Router struct {
	*httprouter.Router
}

func NewDefaultRouter(c *Config) *Router {
	r := httprouter.New()

	return New(r, c)
}

func New(r *httprouter.Router, c *Config) *Router {
	router := &Router{r}
	router.Apply(c)

	return router
}

func (r *Router) Apply(c *Config) {
	r.applyRoutes(c)
	r.applyGroups(c)
}

func (r *Router) applyRoutes(c *Config) {
	for _, route := range c.Routes {
		h := wrap(route.Handler, route.Pipeline)
		h = wrap(h, c.Pipeline)

		r.Handle(route.Method, c.BasePath+route.Path, h)
	}
}

func (r *Router) applyGroups(c *Config) {
	for _, group := range c.Groups {
		group.apply(r, c)
	}
}
