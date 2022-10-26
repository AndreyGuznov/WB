package dto

type ForecastDTO struct {
	Name           string  `json:"name,omitempty"`
	Country        string  `json:"country,omitempty"`
	TempAvg        float64 `json:"tempAvg,omitempty"`
	AvailableDates []int64 `json:"availableDates,omitempty"`
}
