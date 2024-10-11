package models

type GeoData struct {
	Name        string  `json:"name"`
	LongLatData [][]int `json:"longlat"`
	Email       string  `json:"email,omitempty"`
}
