package locations

import (
	"net/http"
	"strconv"
	"sync"
	"weather/tool/log"
	"weather/tool/server"
	"weather/tool/service"

	"github.com/gorilla/mux"
)

type LocationController struct{}

var (
	instance *LocationController
	once     sync.Once
)

func GetInstance() *LocationController {
	once.Do(func() {
		instance = &LocationController{}
		log.Debug("Initialized locations controller")
	})

	return instance
}

// GetRoutes returns the list of all the routes that this handler exposes
func (lc *LocationController) GetRoutes() []server.Route {
	routes := make([]server.Route, 0)

	routes = append(routes, server.Route{Name: "Get All Locations", Method: http.MethodGet, Pattern: "/locations",
		HandlerFunc: lc.getAllLocations})

	routes = append(routes, server.Route{Name: "Get detailed location information at specified time", Method: http.MethodGet, Pattern: "/locations/{locationId}",
		HandlerFunc: lc.getLocationForecastDetails})

	return routes
}

func (lc *LocationController) getAllLocations(w http.ResponseWriter, _ *http.Request) {
	locationDtos, err := service.GetAllLocations()
	if err != nil {
		server.WriteResponse(w, http.StatusInternalServerError, err)
	}
	server.WriteResponse(w, http.StatusOK, locationDtos)
}

func (lc *LocationController) getLocationForecastDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationId, err := strconv.Atoi(vars["locationId"])
	t := r.URL.Query().Get("time")
	if len(t) == 0 {
		server.WriteResponse(w, http.StatusBadRequest, server.NewError(server.BadRequestError, "Missing time"))
	}

	timeReq, err := strconv.ParseInt(t, 10, 0)

	if err != nil {
		server.WriteResponse(w, http.StatusInternalServerError, err)
	}

	details, err := service.GetForecastDetails(locationId, timeReq)

	server.WriteResponse(w, http.StatusOK, details)
}
