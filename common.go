package httprouter_middleware

import "github.com/julienschmidt/httprouter"

func wrap(handler httprouter.Handle, pipeline Pipeline) httprouter.Handle {
	if len(pipeline) == 0 {
		return handler
	}

	wrappedHandler := handler

	for i := len(pipeline) - 1; i >= 0; i-- {
		wrappedHandler = pipeline[i](wrappedHandler)
	}

	return wrappedHandler
}
