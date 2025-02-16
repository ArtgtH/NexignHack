package structs

import (
	"backend/src/service/messages"
	"github.com/google/uuid"
)

type FileTaskResponse struct {
	ID       uuid.UUID                `json:"id"`
	Type     string                   `json:"type"`
	Messages []messages.MessageResult `json:"messages"`
}

type TextTaskRequest struct {
	Text string `json:"text"`
}

type TextTaskResponse struct {
	Text   string `json:"text"`
	Result int    `json:"result"`
}
