package service

import (
	"database/sql"
	"fmt"
	"weather/tool/db"
	"weather/tool/dto"
	"weather/tool/integrations"
	"weather/tool/log"
)

func GetForecastByLocationId(locId int) (*dto.ForecastDTO, error) {
	loc, err := db.GetLocationById(locId)
	if err != nil {
		return nil, err
	}

	fcs, err := db.FindForecastByLocationId(locId)
	if err != nil {
		return nil, err
	}

	if len(fcs) == 0 {
		log.Info(fmt.Sprintf("No forecast for %d", locId))
	}

	dates := make([]int64, 0)

	var tempAvg float64

	for _, fc := range fcs {
		dates = append(dates, fc.Timestamp)
		tempAvg += fc.Temp
	}

	tempAvg = tempAvg / float64(len(fcs))

	dto := dto.ForecastDTO{
		Name:           loc.Name,
		Country:        loc.Country,
		TempAvg:        tempAvg,
		AvailableDates: dates,
	}

	return &dto, nil
}

func GetAllLocations() ([]*dto.LocationDTO, error) {
	locs, err := db.FindAllLocations()
	if err != nil {
		return nil, err
	}

	locationDtos := make([]*dto.LocationDTO, 0)
	for _, loc := range locs {
		dto := dto.LocationDTO{
			Id:      loc.Id,
			City:    loc.Name,
			Country: loc.Country,
		}

		locationDtos = append(locationDtos, &dto)
	}

	return locationDtos, nil
}

func GetForecastDetails(locationId int, timestamp int64) (*dto.LocationDTO, error) {
	loc, err := db.GetLocationById(locationId)
	if err != nil {
		return nil, err
	}

	fc, err := db.GetDetailForecast(locationId, timestamp)
	if err != nil {
		return nil, err
	}

	return &dto.LocationDTO{
		Id:       loc.Id,
		City:     loc.Name,
		Country:  loc.Country,
		Forecast: fc,
	}, nil

}

func InsertLocation(name string) error {
	loc, err := integrations.GetLocation(name)
	if err != nil {
		return err
	}

	id, err := db.GetLocationByName(name)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if err == sql.ErrNoRows {
		id, err = db.InsertOrUpdateLocation(loc)
		if err != nil {
			return err
		}
	}

	forc, err := integrations.GetForecast(loc.Lat, loc.Lng)
	if err != nil {
		return err
	}

	for i := range forc {
		forc[i].LocationId = id
	}

	db.InsertOrUpdateForecast(forc)
	if err != nil {
		return err
	}
	return nil
}
