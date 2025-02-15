package structs

import (
	"backend/src/service/structs"
	"github.com/google/uuid"
)

type JSONTaskRequest struct {
	Data interface{} `json:"data" required:"true"`
}

type CreatedFullTask struct {
	ID       uuid.UUID         `json:"id"`
	Type     string            `json:"type"`
	Messages []structs.Message `json:"messages"`
}

type ResultTask struct {
	ID       uuid.UUID               `json:"id"`
	Type     string                  `json:"type"`
	Messages []structs.MessageResult `json:"messages"`
}

type TextTaskRequest struct {
	Text string `json:"text"`
}

type TextTaskResponse struct {
	Text   string `json:"text"`
	Result int    `json:"result"`
}
