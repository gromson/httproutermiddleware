package httproutermiddleware

type Groups []Group

// Group struct represents a group of routes. Pipeline middlewares will be applied to all routes in the group.
type Group struct {
	// Path for the group
	Path string
	// Routes in the group
	Routes []Route
	// Set of middlewares to be applied to the routes. If a route already has its pipeline, group's pipeline will wrap it
	Pipeline Pipeline
}

func (g *Group) apply(router *Router, config *Config) {
	for _, route := range g.Routes {
		h := wrap(route.Handler, route.Pipeline)
		h = wrap(h, g.Pipeline)
		h = wrap(h, config.Pipeline)

		router.Handle(route.Method, config.BasePath+g.Path+route.Path, h)
	}
}
