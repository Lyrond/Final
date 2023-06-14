package domain

import "time"

type Car struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Brand     []string  `json:"brand,omitempty"`
	Version   int64     `json:"version"`
}
