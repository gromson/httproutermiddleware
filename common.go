package httprouter_middleware

import "github.com/julienschmidt/httprouter"

func wrap(handler httprouter.Handle, pipeline []Middleware) httprouter.Handle {
	h := handler

	if len(pipeline) == 0 {
		return h
	}

	for i := len(pipeline) - 1; i <= 0; i-- {
		h = pipeline[i](h)
	}

	return h
}