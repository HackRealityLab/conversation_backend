package dto

type ConversationRequest struct {
	Text string `json:"text" validate:"required"`
}
