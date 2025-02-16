package messages

import "github.com/google/uuid"

type Message struct {
	UserID      string `json:"userID"`
	SubmitDate  string `json:"submitDate"`
	MessageText string `json:"messageText"`
}

type MessageResult struct {
	UserID      string `json:"userID"`
	SubmitDate  string `json:"submitDate"`
	MessageText string `json:"messageText"`
	Result      int    `json:"result"`
}

type CreatedFullTask struct {
	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	Messages []Message `json:"messages"`
}
