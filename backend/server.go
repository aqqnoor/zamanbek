package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"unicode/utf8"

	"zmb-assistant/services"

	"golang.org/x/text/encoding/charmap"
)

var sessions = NewSessionStore()

func startServer() {
	mux := http.NewServeMux()

	// Хелпер: всегда JSON UTF-8
	jsonUTF8 := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/llm/ping", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		llm := services.NewLLM(os.Getenv("OPENAI_BASE_URL"), os.Getenv("OPENAI_API_KEY"))
		reply, err := llm.Chat("gpt-4o-mini",
			[]services.ChatMessage{{Role: "user", Content: "Скажи: pong"}}, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"reply": reply})
	})

	mux.HandleFunc("/session/new", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		id := sessions.EnsureSessionID(w, r)
		sessions.Reset(id)
		_ = json.NewEncoder(w).Encode(map[string]any{"session_id": id, "status": "new"})
	})

	mux.HandleFunc("/session/clear", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		id := sessions.EnsureSessionID(w, r)
		sessions.Reset(id)
		_ = json.NewEncoder(w).Encode(map[string]any{"session_id": id, "status": "cleared"})
	})

	mux.HandleFunc("/session/history", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		id := sessions.EnsureSessionID(w, r)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"session_id": id,
			"history":    sessions.Get(id),
		})
	})

	mux.HandleFunc("/llm/chat", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		if r.Method != http.MethodPost {
			http.Error(w, "use POST", http.StatusMethodNotAllowed)
			return
		}
		id := sessions.EnsureSessionID(w, r)

		var body struct {
			Message string  `json:"message"`
			Temp    float32 `json:"temperature,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// Если сообщение невалидное UTF-8 (часто из Windows), пробуем CP1251 → UTF-8
		if !utf8.ValidString(body.Message) {
			if dec, err := charmap.Windows1251.NewDecoder().String(body.Message); err == nil {
				body.Message = dec
			}
		}

		if body.Temp == 0 {
			body.Temp = 0.2
		}

		history := sessions.Get(id)
		if len(history) == 0 {
			history = append(history, services.ChatMessage{
				Role: "system",
				Content: "Сен — ZamanBank ассистентісің. Негізгі тіл: қазақ тілі. " +
					"Жауаптарың қысқа, жылы әрі нақты болсын. " +
					"Исламдық қаржыландыру қағидаларын сақта: риба/өсім жоқ, " +
					"спекуляциядан аулақ, активке негізделген ұсыныстар. " +
					"Қажет болса орысша түсіндір, бірақ әуелі қазақша жауап бер.",
			})
		}

		userMsg := services.ChatMessage{Role: "user", Content: body.Message}
		sessions.Append(id, userMsg)
		history = append(history, userMsg)

		llm := services.NewLLM(os.Getenv("OPENAI_BASE_URL"), os.Getenv("OPENAI_API_KEY"))
		reply, err := llm.Chat("gpt-4o-mini", history, body.Temp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		sessions.Append(id, services.ChatMessage{Role: "assistant", Content: reply})

		_ = json.NewEncoder(w).Encode(map[string]any{
			"session_id":  id,
			"reply":       reply,
			"history_len": len(sessions.Get(id)),
		})
	})

	// Раздаём статические файлы из frontend/
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../"))))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
