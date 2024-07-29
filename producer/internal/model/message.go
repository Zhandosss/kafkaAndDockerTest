package model

import "time"

type Message struct {
	ID          string    `json:"id" db:"id"`
	Content     string    `json:"content" db:"content"`
	IsProcessed bool      `json:"is_processed" db:"is_processed"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}
