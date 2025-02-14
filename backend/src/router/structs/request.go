package structs

import "backend/src/service/converter"

type JSONTaskRequest struct {
	Data interface{} `json:"data" required:"true"`
}

type CreatedTask struct {
	ID       uint                `json:"id"`
	Type     string              `json:"type"`
	Messages []converter.Message `json:"messages"`
}
