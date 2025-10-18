package services

// ChatMessage represents a chat message exchanged with the language model.
type ChatMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}