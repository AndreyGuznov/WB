package integrations

import (
	"weather/tool/db"
	"weather/tool/log"
)

func RefreshData() {

	log.Info("Refreshing data")

	ch := make(chan []*db.Forecast, 1)

	go func() {
		locations, err := db.FindAllLocations()
		if err != nil {
			log.Err("Error of Refresh", err)
		}
		for i := range locations {
			forc, err := GetForecast(locations[i].Lat, locations[i].Lng)
			if err != nil {
				log.Err("Error of Refresh", err)
			}
			for j := range forc {
				forc[j].LocationId = locations[i].Id
			}
			ch <- forc
		}
	}()

	go db.InsertOrUpdateForecast(<-ch)

}
