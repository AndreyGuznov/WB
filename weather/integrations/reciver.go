package integrations

import "weather/tool/db"

func RefreshData() error {
	forecastList := []*db.Forecast{}
	locations, err := db.FindAllLocations()
	if err != nil {
		return err
	}
	for i := range locations {
		forc, err := GetForecast(locations[i].Lat, locations[i].Lng)
		if err != nil {
			return err
		}
		for j := range forc {
			forc[j].LocationId = locations[i].Id
		}
		forecastList = append(forecastList, forc...)
	}
	err = db.InsertOrUpdateForecast(forecastList)
	if err != nil {
		return err
	}

	return nil
}
