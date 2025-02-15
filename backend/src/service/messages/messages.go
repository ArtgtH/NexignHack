package messages

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
