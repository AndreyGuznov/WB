package forecasts

import (
	"net/http"
	"strconv"
	"sync"
	"weather/tool/log"
	"weather/tool/server"
	"weather/tool/service"

	"github.com/gorilla/mux"
)

type ForecastController struct{}

var (
	instance *ForecastController
	once     sync.Once
)

func GetInstance() *ForecastController {
	once.Do(func() {
		instance = &ForecastController{}
		log.Debug("Initialized locations controller")
	})

	return instance
}

// GetRoutes returns the list of all the routes that this handler exposes
func (fc *ForecastController) GetRoutes() []server.Route {
	routes := make([]server.Route, 0)

	routes = append(routes, server.Route{Name: "Get location forecast", Method: http.MethodGet, Pattern: "/forecast/{locationId}",
		HandlerFunc: fc.getForecast})

	return routes
}

func (fc *ForecastController) getForecast(w http.ResponseWriter, r *http.Request) {
	// Read the location id from the url path.
	vars := mux.Vars(r)
	locationId, err := strconv.ParseInt(vars["locationId"], 10, 0)
	if err != nil {
		server.WriteResponse(w, http.StatusBadRequest, server.NewError(server.BadRequestError, "Invalid location id"))
	}

	forecast, err := service.GetForecastByLocationId(int(locationId))
	if err != nil {
		server.WriteResponse(w, http.StatusInternalServerError, err)
	}

	server.WriteResponse(w, http.StatusOK, forecast)
}
