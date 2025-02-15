package structs

import (
	"backend/src/service/messages"
	"github.com/google/uuid"
)

type JSONTaskRequest struct {
	Data interface{} `json:"data" required:"true"`
}

type CreatedFullTask struct {
	ID       uuid.UUID          `json:"id"`
	Type     string             `json:"type"`
	Messages []messages.Message `json:"messages"`
}

type ResultTask struct {
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
