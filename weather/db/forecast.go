package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"weather/tool/log"
)

// Forecast DB type
type Forecast struct {
	LocationId int     `db:"location_id" json:"locationId,omitempty"`
	Timestamp  int64   `db:"timestamp" json:"dt,omitempty"`
	Temp       float64 ` db:"temp" json:"temp,omitempty"`
	Data       string  `db:"data" json:"data,omitempty"`
}

func FindForecastByLocationId(locationId int) ([]*Forecast, error) {

	now := time.Now()

	log.Info("Trying to get forcastById from Forecasts")

	query := fmt.Sprintf(`SELECT
			location_id,
			timestamp,
			temp
		FROM %s WHERE location_id=$1 AND timestamp>=$2`, ForecastTable)

	rows, err := GetConn().Instance.Queryx(query, locationId, now.Unix())

	if err != nil {
		return nil, err
	}

	forecast := []*Forecast{}

	for rows.Next() {
		fc := Forecast{}

		if err := rows.StructScan(&fc); err != nil {
			return nil, err
		}

		forecast = append(forecast, &fc)
	}
	log.Info("Success of geting forcastById from Forecasts")
	return forecast, nil
}

func GetDetailForecast(locationId int, timestamp int64) (*Forecast, error) {
	var forecast Forecast

	log.Info("Trying to get DetailForecast info from ForecastTable")

	query := fmt.Sprintf(`SELECT location_id, timestamp, data FROM %s WHERE location_id=$1 and timestamp=$2`, ForecastTable)
	if err := GetConn().Instance.Get(&forecast, query, locationId, timestamp); err != nil {
		return nil, err
	}
	log.Info("Success of geting DetailForecast from ForecastTable")
	return &forecast, nil
}

func getNext(count *int) *int {
	*count += 1
	return count
}

func InsertOrUpdateForecast(forecasts []*Forecast) {
	paramCount := 0
	valsPlaceholders := make([]string, 0, len(forecasts))
	vals := make([]interface{}, 0, 4*len(forecasts))

	for _, fc := range forecasts {
		bindVars := "($" + strconv.Itoa(*getNext(&paramCount)) + ",$" + strconv.Itoa(*getNext(&paramCount)) + ",$" + strconv.Itoa(*getNext(&paramCount)) + ",$" + strconv.Itoa(*getNext(&paramCount)) + ")"

		valsPlaceholders = append(valsPlaceholders, bindVars)
		vals = append(vals, fc.LocationId, fc.Timestamp, fc.Temp, fc.Data)
	}

	query := fmt.Sprintf(`INSERT INTO %s (
		location_id,
		timestamp,                                                                   
		temp,
		data
	  )
	  VALUES %s
	  ON CONFLICT (location_id, timestamp) DO UPDATE
	  SET temp=EXCLUDED.temp, data=EXCLUDED.data
	  `, ForecastTable, strings.Join(valsPlaceholders, ", "))

	if _, err := GetConn().Instance.Exec(query, vals...); err != nil {
		log.Err("Some thing wrong with InsertOrUpdateForecast", err)
	}
}
