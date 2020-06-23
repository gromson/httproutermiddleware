package httprouter_middleware

// Group struct represents a group of routes. Pipeline middlewares will be applied to all routes in the group
type Group struct {
	// Routes in the group
	Routes []Route
	// Set of middlewares to be applied to the routes. If a route already has its pipeline, group's pipeline will wrap it
	Pipeline []Middleware
}

func (g *Group) apply(r *Router) {
	for _, route := range g.Routes {
		h := wrap(route.wrap(), g.Pipeline)
		r.Handle(route.Method, route.Path, h)
	}
}