package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"zmb-assistant/jobs"
)

func TrainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed); return
	}
	base := os.Getenv("OPENAI_BASE_URL")
	if base == "" { base = "https://openai-hub.neuraldeep.tech" }
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" { http.Error(w, "OPENAI_API_KEY missing", 500); return }

	if err := jobs.TrainIndex(base, key); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway); return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"status":"ok"})
}
