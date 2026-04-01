package model

import "time"

type URL struct {
	ID        int64
	Original  string
	Code      string
	CreatedAt time.Time
}
