package db

import (
	"fmt"
	"weather/tool/log"
)

type Location struct {
	Id      int     `db:"id"`
	Name    string  `db:"name" json:"name,omitempty"`
	Country string  `db:"country" json:"country,omitempty"`
	Lat     float64 `db:"lat" json:"lat,omitempty"`
	Lng     float64 `db:"lng" json:"lon,omitempty"`
}

func FindAllLocations() ([]*Location, error) {
	locations := []*Location{}

	log.Info("Trying to get all info from LocationsTable")

	query := fmt.Sprintf(`SELECT 
			id, 
			name,
			country
		FROM %s
		ORDER BY name`, LocationsTable)

	acRows, err := GetConn().Instance.Queryx(query)

	if err != nil {
		return nil, err
	}

	for acRows.Next() {
		location := Location{}

		if err := acRows.StructScan(&location); err != nil {
			return nil, err
		}

		locations = append(locations, &location)
	}

	log.Info("Success of get all info from LocationsTable")

	return locations, nil
}

func GetLocationById(id int) (*Location, error) {
	var location Location

	log.Info("Trying to get LocationById info from LocationsTable")

	query := fmt.Sprintf(`SELECT 
			id, 
			name,
			country,
			lat,
			lng
		FROM %s
		WHERE id=$1`, LocationsTable)
	if err := GetConn().Instance.Get(&location, query, id); err != nil {
		return nil, err
	}

	log.Info("Success of geting LocationById from LocationsTable")

	return &location, nil
}

func InsertOrUpdateLocation(location *Location) (int, error) {

	log.Info("Trying to isert info in LocationsTable")
	id := 0
	query := fmt.Sprintf(`INSERT INTO %s (
		name,
		country,
		lat,
		lng
	  )  
	  VALUES ($1, $2, $3, $4)  
	  ON CONFLICT (name,country) DO NOTHING
	  RETURNING id
	  `, LocationsTable)

	err := GetConn().Instance.QueryRow(query, location.Name, location.Country, location.Lat, location.Lng).Scan(&id)
	if err != nil {
		return 0, err
	}

	log.Info("Succes of iserting in LocationsTable")

	return id, nil

}

// TODO err with same "names"
func GetLocationByName(name string) (*Location, error) {
	var location Location

	query := fmt.Sprintf(`SELECT 
		id, 
		name,
		country,
		lat,
		lng
	  FROM %s
	  WHERE name=$1`, LocationsTable)
	if err := GetConn().Instance.Get(&location, query, name); err != nil {
		return nil, err
	}

	return &location, nil
}
