package dto

import (
	"weather/tool/db"
)

type LocationDTO struct {
	Id       int          `json:"id,omitempty"`
	City     string       `json:"city,omitempty"`
	Country  string       `json:"country,omitempty"`
	Forecast *db.Forecast `json:"details,omitempty"`
}
