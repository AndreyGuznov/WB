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
	geoSvcUrl    = config.Get("GEO_URL")
	geoSvcApiKey = config.Get("GEO_API_KEY")
)

func GetLocation(name string) (*db.Location, error) {
	geoUrl := geoSvcUrl + name + "&limit=1&appid=" + geoSvcApiKey
	log.Debug("Calling geocoding-api get location = " + geoUrl)

	// Call inventory service REST API
	resp, err := http.Get(geoUrl)
	if err != nil {
		log.Err("Unable to read location details from geocoding-api", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Info("No location found for given name = " + name)
		return nil, errors.New("Location not found in geocoding-api")
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("%d status code from geocoding-api service", resp.StatusCode))
		log.Err("Unable to read location details from geocoding-api", err)
		return nil, err
	}

	// Unmarshal the JSON returned by the API
	decoder := json.NewDecoder(resp.Body)
	var locs []*db.Location
	err = decoder.Decode(&locs)
	if err != nil {
		log.Err(fmt.Sprintf("Unable to unmarshal location %s details returned from geocoding-api", name), err)
		return nil, err
	}

	if len(locs) == 0 {
		log.Err(fmt.Sprintf("Unable to get location from API %s after unmarshal", name), err)
		return nil, err
	}

	loc := locs[0]

	return loc, nil
}
