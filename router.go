package httprouter_middleware

import (
	"github.com/julienschmidt/httprouter"
)

// Middleware wraps the next handler
type Middleware func(next httprouter.Handle) httprouter.Handle

// Pipeline the pipeline of middlawares
type Pipeline []Middleware

type Config struct {
	// Groups
	Groups Groups
	// Routes
	Routes Routes
}

// Router wrapper for httprouter.Router
type Router struct {
	*httprouter.Router
}

func NewDefaultRouter(c *Config) *Router {
	r := httprouter.New()
	router := &Router{r}
	router.Apply(c)

	return router
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
		h := route.wrap()
		r.Handle(route.Method, route.Path, h)
	}
}

func (r *Router) applyGroups(c *Config) {
	for _, group := range c.Groups {
		group.apply(r)
	}
}
