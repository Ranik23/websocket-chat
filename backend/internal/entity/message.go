package entity

type Message struct {
	SenderName string `json:"sender"`
	Text       string `json:"text"`
	Room       string `json:"room"`
}
