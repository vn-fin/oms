package models

// EnsembleWeight represents the ensemble id
type EnsembleWeight struct {
	EnsembleID string  `json:"ensemble_id"`
	Weight     float64 `json:"weight"`
}
