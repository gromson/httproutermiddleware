package httproutermiddleware_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mw "github.com/gromson/httproutermiddleware"
	"github.com/julienschmidt/httprouter"
)

type cxtKey int

const ctxTestKey cxtKey = iota

func TestRouter_Apply(t *testing.T) {
	t.Parallel()
	t.Run("standalone route", testStandaloneRoute)
	t.Run("grouped route", testGroupedRoute)
	t.Run("no pipeline", testNoPipeline)
}

func testStandaloneRoute(t *testing.T) {
	// Arrange
	t.Parallel()

	config := &mw.Config{
		Routes: mw.Routes{
			mw.Route{
				Path:    "/route",
				Handler: handler,
				Method:  "GET",
				Pipeline: mw.Pipeline{
					routeMiddleware,
					routeMiddleware2,
				},
			},
		},
		Pipeline: mw.Pipeline{
			globalMiddleware,
			globalMiddleware2,
		},
	}

	// Act & Assert
	executeAndAssert(t, config, "/route", "global global2 route route2")
}

func testGroupedRoute(t *testing.T) {
	// Arrange
	t.Parallel()

	config := &mw.Config{
		Groups: mw.Groups{
			mw.Group{
				Routes: mw.Routes{
					mw.Route{
						Path:    "/route",
						Handler: handler,
						Method:  "GET",
						Pipeline: mw.Pipeline{
							routeMiddleware,
							routeMiddleware2,
						},
					},
				},
				Pipeline: mw.Pipeline{
					groupMiddleware,
					groupMiddleware2,
				},
			},
		},
		Pipeline: mw.Pipeline{
			globalMiddleware,
			globalMiddleware2,
		},
	}

	// Act & Assert
	executeAndAssert(t, config, "/route", "global global2 group group2 route route2")
}

func testNoPipeline(t *testing.T) {
	// Arrange
	t.Parallel()

	config := &mw.Config{
		Routes: mw.Routes{
			mw.Route{
				Path:    "/route",
				Handler: handler,
				Method:  "GET",
			},
		},
	}

	// Act & Assert
	executeAndAssert(t, config, "/route", "null")
}

func executeAndAssert(t *testing.T, config *mw.Config, route, expectedValue string) {
	// Arrange
	t.Helper()

	router := mw.NewDefaultRouter(config)

	srv := httptest.NewServer(router)
	defer srv.Close()

	// Act
	res, err := http.Get(srv.URL + route)
	if err != nil {
		t.Fatal("Could not get a response from the route")
	}

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Could not read a response body")
	}

	_ = res.Body.Close()

	// Assert
	if string(bodyData) != expectedValue {
		t.Fatalf(`Expected response body from /route: "%s", "%s" given`, expectedValue, bodyData)
	}
}

func handler(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	value := req.Context().Value(ctxTestKey)

	if value == nil {
		value = "null"
	}

	if _, err := writer.Write([]byte(value.(string))); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
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

func getRequestWithUpdatedContext(req *http.Request, val string) *http.Request {
	newVal := val
	ctx := req.Context()

	ctxVal := ctx.Value(ctxTestKey)
	if ctxVal != nil {
		ctxValStr, _ := ctxVal.(string)
		newVal = ctxValStr + " " + val
	}

	return req.WithContext(context.WithValue(ctx, ctxTestKey, newVal))
}
