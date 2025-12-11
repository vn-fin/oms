package models

// TimeSeriesData represents the historical data for a strategy with times and values
type TimeSeriesData struct {
	Time  []float64 `json:"time" example:"1625097600.0,1625184000.0"`
	Value []float64 `json:"value" example:"1500.5,1600"`
}
