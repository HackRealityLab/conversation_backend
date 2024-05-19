package dto

import "time"

type ConversationRequest struct {
	Text string `json:"text" validate:"required"`
}

type Record struct {
	ID          int       `json:"id"`
	Text        string    `json:"text"`
	AudioName   string    `json:"audio_name"`
	CreatedAt   time.Time `json:"created_at"`
	GoodPercent int       `json:"good_percent"`
	BadPercent  int       `json:"bad_percent"`
}
