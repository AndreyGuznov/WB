package integrations

import (
	"sync"
	"weather/tool/db"
	"weather/tool/log"
)

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	fmt.Printf("%s took %s", name, elapsed)
// }

func RefreshData() error {

	// defer timeTrack(time.Now(), "Fetching refresh")

	log.Info("Refreshing data")

	done := make(chan struct{}, 1)
	var waitgroup sync.WaitGroup

	forecastList := []*db.Forecast{}

	locations, err := db.FindAllLocations()

	if err != nil {
		return err
	}

	waitgroup.Add(len(locations))

	go func() chan struct{} {

		for i := range locations {

			waitgroup.Add(1)

			forc, err := GetForecast(locations[i].Lat, locations[i].Lng)
			if err != nil {
				log.Err("Error of Refresh", err)
			}

			for j := range forc {
				forc[j].LocationId = locations[i].Id
			}

			forecastList = append(forecastList, forc...)
		}

		waitgroup.Done()

		done <- struct{}{}

		return done
	}()

	go func() error {

		<-done

		err = db.InsertOrUpdateForecast(forecastList)

		if err != nil {
			return err
		}

		return err

	}()

	return nil
}
