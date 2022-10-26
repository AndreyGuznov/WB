package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"weather/tool/log"
)

// Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(http.ResponseWriter, *http.Request)
}

// Controller interface
type Controller interface {
	GetRoutes() []Route
	// Cleanup()
}

func WriteResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			log.Err("Could not marshal task to json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error in marshalling"))
			return
		}
		w.WriteHeader(statusCode)
		w.Write(jsonBody)
	}
}

const contextRoot string = "/weather"

// NewHTTPHandler creates the handler for the various rest routes
func NewHTTPHandler(routesHandler ...Controller) http.Handler {
	muxHandler := mux.NewRouter().StrictSlash(true)
	for _, routeHandler := range routesHandler {
		routes := routeHandler.GetRoutes()
		for _, route := range routes {
			muxHandler.PathPrefix(contextRoot).Methods(route.Method).Path(route.Pattern).
				Name(route.Name).HandlerFunc(route.HandlerFunc)
		}
	}
	return muxHandler
}
