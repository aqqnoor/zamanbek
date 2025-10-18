package services

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type FeedbackEntry struct {
	Time     string `json:"time"`
	Query    string `json:"query"`
	Reply    string `json:"reply"`
	Expected string `json:"expected,omitempty"`
}

const feedbackPath = "data/feedback_log.json"

// LogFeedback — сохраняет ответы для разборов (минимальный критерий «слабости» убран — логируем все)
func LogFeedback(q, reply string) {
	entry := FeedbackEntry{
		Time:  time.Now().Format(time.RFC3339),
		Query: q,
		Reply: reply,
	}
	f, _ := os.OpenFile(feedbackPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	data, _ := json.Marshal(entry)
	f.Write(data)
	f.Write([]byte("\n"))
	log.Printf("[selftrain] logged reply for: %s (len=%d)", truncate(q, 60), len(reply))
}

func truncate(s string, n int) string {
	rs := []rune(s)
	if len(rs) <= n {
		return s
	}
	return string(rs[:n]) + "…"
}
