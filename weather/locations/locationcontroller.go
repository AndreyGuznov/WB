package locations

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"weather/tool/integrations"
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

	routes = append(routes, server.Route{Name: "Add Location", Method: http.MethodPost, Pattern: "/locations",
		HandlerFunc: lc.addLocations})
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

func (lc *LocationController) addLocations(w http.ResponseWriter, r *http.Request) {
	reqBody := integrations.JSONMap{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.WriteResponse(w, http.StatusBadRequest, err)
	}

	name, ok := reqBody["name"].(string)
	if !ok {
		server.WriteResponse(w, http.StatusBadRequest, errors.New("failed to get name"))
	}

	err := service.InsertLocation(name)
	if err != nil {
		server.WriteResponse(w, http.StatusInternalServerError, err)
	}
	server.WriteResponse(w, http.StatusOK, "New location Added")
}
