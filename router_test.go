package httprouter_middleware

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

const contextValueKey = "test"

func TestRouter_Apply(t *testing.T) {
	t.Parallel()
	t.Run("standalone route", testStandaloneRoute)
	t.Run("grouped route", testGroupedRoute)
	t.Run("no pipeline", testNoPipeline)
}

func testStandaloneRoute(t *testing.T) {
	// Arrange
	config := &Config{
		Routes: Routes{
			Route{
				Path:    "/route",
				Handler: handler,
				Method:  "GET",
				Pipeline: Pipeline{
					routeMiddleware,
					routeMiddleware2,
				},
			},
		},
		Pipeline: Pipeline{
			globalMiddleware,
			globalMiddleware2,
		},
	}

	// Act & Assert
	executeAndAssert(t, config, "/route", "global global2 route route2")
}

func testGroupedRoute(t *testing.T) {
	// Arrange
	config := &Config{
		Groups: Groups{
			Group{
				Routes: Routes{
					Route{
						Path:    "/route",
						Handler: handler,
						Method:  "GET",
						Pipeline: Pipeline{
							routeMiddleware,
							routeMiddleware2,
						},
					},
				},
				Pipeline: Pipeline{
					groupMiddleware,
					groupMiddleware2,
				},
			},
		},
		Pipeline: Pipeline{
			globalMiddleware,
			globalMiddleware2,
		},
	}

	// Act & Assert
	executeAndAssert(t, config, "/route", "global global2 group group2 route route2")
}

func testNoPipeline(t *testing.T) {
	// Arrange
	config := &Config{
		Routes: Routes{
			Route{
				Path:    "/route",
				Handler: handler,
				Method:  "GET",
			},
		},
	}

	// Act & Assert
	executeAndAssert(t, config, "/route", "null")
}

func executeAndAssert(t *testing.T, config *Config, route, expectedValue string) {
	// Arrange
	router := NewDefaultRouter(config)
	srv := httptest.NewServer(router)
	defer srv.Close()

	// Act
	res, err := http.Get(srv.URL + route)
	if err != nil {
		t.Fatal("Could not get a response from the route")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Could not read a response body")
	}
	_ = res.Body.Close()

	// Assert
	if string(b) != expectedValue {
		t.Fatalf(`Expected response body from /route: "%s", "%s" given`, expectedValue, b)
	}
}

func handler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	value := r.Context().Value(contextValueKey)

	if value == nil {
		value = "null"
	}

	if _, err := w.Write([]byte(value.(string))); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func globalMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, getRequestWithUpdatedContext(r, "global"), p)
	}
}

func globalMiddleware2(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, getRequestWithUpdatedContext(r, "global2"), p)
	}
}

func groupMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, getRequestWithUpdatedContext(r, "group"), p)
	}
}

func groupMiddleware2(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, getRequestWithUpdatedContext(r, "group2"), p)
	}
}

func routeMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, getRequestWithUpdatedContext(r, "route"), p)
	}
}

func routeMiddleware2(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		next(w, getRequestWithUpdatedContext(r, "route2"), p)
	}
}

func getRequestWithUpdatedContext(r *http.Request, val string) *http.Request {
	newVal := val
	ctx := r.Context()
	ctxVal := ctx.Value(contextValueKey)
	if ctxVal != nil {
		newVal = ctxVal.(string) + " " + val
	}

	return r.WithContext(context.WithValue(ctx, contextValueKey, newVal))
}
