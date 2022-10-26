package main

import (
	"errors"
	"net/http"
	"os"
	"time"
	"weather/tool/config"
	"weather/tool/db"
	"weather/tool/forecasts"
	"weather/tool/integrations"
	"weather/tool/locations"
	"weather/tool/log"
	"weather/tool/server"
)

var (
	attempts = config.GetInt("CONN_ATTEMPTS")
)

func main() {
	cityList := []string{"Barcelona", "Boston", "London", "Moscow", "Yerevan"}
	for i := 0; i < attempts; i++ {
		if db.GetConn() != nil {
			break
		}
		log.Err("Unable to connect to Postgres. Retrying...", nil)
		time.Sleep(5 * time.Second)
	}

	if db.GetConn() == nil {
		log.Err("Failed to connect to Postgres", errors.New("broken pipe"))
		os.Exit(1)
	}

	// // fmt.Println(db.FindForecastByLocationId(1))
	// fmt.Println(time.Unix(1666753200, 0))
	// // fmt.Printf("%T", time.Unix(1666753200, 0))
	// fmt.Println(db.GetDetailForecast(1, 1666807200))
	// return

	err := integrations.RefreshData(cityList)
	if err != nil {
		log.Err("", err)
	}

	handler := getHTTPHandler()

	httpServer := server.NewHTTPRestServer(":"+config.Get("PORT"), handler)
	_ = httpServer.Serve()

}

func getControllers() []server.Controller {
	controllers := make([]server.Controller, 0)
	//Add all the controllers here
	controllers = append(controllers, locations.GetInstance())
	controllers = append(controllers, forecasts.GetInstance())

	return controllers
}

// Creates the handler for the various rest routes
func getHTTPHandler() http.Handler {
	controllers := getControllers()
	return server.NewHTTPHandler(controllers...)
}
