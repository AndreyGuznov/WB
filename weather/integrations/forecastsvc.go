package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"weather/tool/config"
	"weather/tool/db"
	"weather/tool/log"
)

var (
	forecastSvcUrl    = config.Get("FORECAST_URL")
	forecastSvcApiKey = config.Get("FORECAST_API_KEY")
)

type JSONMap = map[string]interface{}

func GetForecast(lat, lon float64) ([]*db.Forecast, error) {
	geoUrl := fmt.Sprintf("http://%slat=%f&lon=%f&appid=%s&units=metric", forecastSvcUrl, lat, lon, forecastSvcApiKey)
	log.Debug("Calling openweathermap service get forecast = " + geoUrl)

	// Call openweathermap service REST API
	resp, err := http.Get(geoUrl)
	if err != nil {
		log.Err("Unable to read location details from forecast service", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Info("No location found for given coordinates = ")
		return nil, errors.New(fmt.Sprintf("%d status code from forecast service", resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("%d status code from forecast service", resp.StatusCode))
		log.Err("Unable to read location details from forecast service", err)
		return nil, err
	}

	result, err := responseToEntity(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func responseToEntity(resp *http.Response) ([]*db.Forecast, error) {
	result := make([]*db.Forecast, 0)

	searchResult := JSONMap{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		log.Err("failed to decode query result from forecast service", err)
		return nil, err
	}

	fcs, ok := searchResult["list"].([]interface{})
	if !ok {
		err := errors.New("failed to decode forecast responses")
		log.Err("failed to decode forecast responses", err)
		return nil, err
	}

	for _, fcJson := range fcs {
		data, _ := json.Marshal(fcJson)

		fc := db.Forecast{
			Data: string(data),
		}

		obj, _ := fcJson.(JSONMap)

		ts, _ := obj["dt"].(float64)

		fc.Timestamp = int64(ts)

		main := obj["main"].(JSONMap)

		fc.Temp = main["temp"].(float64)
		result = append(result, &fc)
	}

	return result, nil
}
