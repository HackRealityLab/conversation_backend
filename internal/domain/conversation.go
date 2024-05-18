package domain

import "time"

type Record struct {
	ID          int
	Text        string
	AudioName   string
	CreatedAt   time.Time
	GoodPercent int
	IsOk        bool
}
