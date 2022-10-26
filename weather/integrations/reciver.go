package integrations

import (
	"weather/tool/db"
)

func RefreshData(cityList []string) error {

	forecastlistdata := []*db.Forecast{}

	for i := range cityList {

		loc, err := GetLocation(cityList[i])
		if err != nil {
			return err
		}

		id, err := db.InsertOrUpdateLocation(loc)
		if err != nil {
			return err
		}

		forc, err := GetForecast(loc.Lat, loc.Lng)
		if err != nil {
			return err
		}

		for i := range forc {
			forc[i].LocationId = id
		}

		forecastlistdata = append(forecastlistdata, forc...)
	}

	err := db.InsertOrUpdateForecast(forecastlistdata)
	if err != nil {
		return err
	}

	return nil
}
