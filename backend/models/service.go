package models

type Service struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	ImageURL        string  `json:"image_url"`
}
