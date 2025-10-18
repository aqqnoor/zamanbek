package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"zmb-assistant/handlers"
	"zmb-assistant/services"

	"golang.org/x/text/encoding/charmap"
)

var (
	sessions             = NewSessionStore()
	systemPromptTemplate = ""
	cfgBaseURL           = ""
	cfgAPIKey            = ""
	cfgBackend           = ""
	cfgModel             = ""
	defaultTemperature   = float32(0.2)

	staticRoot string
)

// --- вспомогательные функции статики ---
func findStaticRoot() string {
	candidates := []string{"..", ".", "../.."}
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates, filepath.Join(exeDir, ".."), exeDir, filepath.Join(exeDir, "../.."))
	}
	seen := map[string]bool{}
	for _, c := range candidates {
		cAbs, _ := filepath.Abs(c)
		if seen[cAbs] {
			continue
		}
		seen[cAbs] = true
		if _, err := os.Stat(filepath.Join(cAbs, "chat.html")); err == nil {
			return cAbs
		}
	}
	wd, _ := os.Getwd()
	return wd
}

func mustLoadSystemPrompt() {
	data, err := ioutil.ReadFile("data/system_prompt.txt")
	if err != nil {
		log.Fatalf("failed to load system prompt: %v", err)
	}
	systemPromptTemplate = string(data)
}
func buildSystemPrompt(lang string) string {
	return strings.ReplaceAll(systemPromptTemplate, "{{LANG}}", lang)
}

// мягкий тримминг истории
const (
	maxHistoryMsgs  = 24
	maxContentChars = 2000
)
func trimHistoryIfNeeded(in []services.ChatMessage) []services.ChatMessage {
	for i := range in {
		if len([]rune(in[i].Content)) > maxContentChars {
			rs := []rune(in[i].Content)
			in[i].Content = string(rs[:maxContentChars])
		}
	}
	if len(in) > maxHistoryMsgs {
		return append([]services.ChatMessage(nil), in[len(in)-maxHistoryMsgs:]...)
	}
	return in
}

func startServer() {
	// конфиг
	cfgBaseURL = os.Getenv("OPENAI_BASE_URL")
	cfgAPIKey = os.Getenv("OPENAI_API_KEY")
	cfgBackend = os.Getenv("LLM_BACKEND") // "bank" or "mock"
	if cfgBackend == "" {
		cfgBackend = "mock"
	}
	cfgModel = os.Getenv("LLM_MODEL")
	if cfgModel == "" {
		cfgModel = "gpt-4o-mini"
	}

	mustLoadSystemPrompt()
	services.LoadProducts()
	services.LoadKnowledge()

	// статика
	staticRoot = findStaticRoot()
	log.Println("[static] serving from:", staticRoot)

	mux := http.NewServeMux()
	jsonUTF8 := func(w http.ResponseWriter) { w.Header().Set("Content-Type", "application/json; charset=utf-8") }

	// --- health ---
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		ready := systemPromptTemplate != ""
		if ready {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ready"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("not ready"))
	})

	// --- admin: reload data ---
	mux.HandleFunc("/admin/reload", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		if r.Method != http.MethodPost {
			http.Error(w, "use POST", http.StatusMethodNotAllowed)
			return
		}
		services.LoadProducts()
		services.LoadKnowledge()
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "reloaded"})
	})

	// --- admin: feedback → train ---
	mux.HandleFunc("/admin/feedback/train", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		file := "data/feedback_log.json"
		data, err := ioutil.ReadFile(file)
		if err != nil {
			http.Error(w, "no feedback log", http.StatusNotFound)
			return
		}
		lines := strings.Split(string(data), "\n")
		added := 0
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var entry services.FeedbackEntry
			if json.Unmarshal([]byte(line), &entry) == nil {
				services.AddFactFromFeedback(entry)
				added++
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"trained_from": added})
	})

	// --- ping ---
	mux.HandleFunc("/llm/ping", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		llm := services.NewLLM(cfgBaseURL, cfgAPIKey, cfgBackend)
		reply, err := llm.Chat(cfgModel, []services.ChatMessage{{Role: "user", Content: "ping"}}, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"reply": reply})
	})

	// --- sessions ---
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
		_ = json.NewEncoder(w).Encode(map[string]any{"session_id": id, "history": sessions.Get(id)})
	})

	// --- debug: detect ---
	mux.HandleFunc("/debug/detect", func(w http.ResponseWriter, r *http.Request) {
		jsonUTF8(w)
		if r.Method != http.MethodPost {
			http.Error(w, "use POST", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Message string `json:"message"`
		}
		bb, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(bb, &body)

		msg := body.Message
		utfOK := utf8.ValidString(msg)
		if !utfOK {
			if dec, err := charmap.Windows1251.NewDecoder().String(msg); err == nil {
				msg = dec
				utfOK = utf8.ValidString(msg)
			}
		}
		runes := []rune(msg)
		limit := 16
		if len(runes) < limit {
			limit = len(runes)
		}
		hexRunes := make([]string, limit)
		for i := 0; i < limit; i++ {
			b := []byte(string(runes[i]))
			hexRunes[i] = hex.EncodeToString(b)
		}
		lang := services.DetectLang(msg)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"raw_json":   string(bb),
			"raw_msg":    body.Message,
			"msg_after":  msg,
			"utf8_valid": utfOK,
			"runes_hex":  hexRunes,
			"lang":       lang,
		})
	})

	// --- chat ---
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

		// кодировка
		if !utf8.ValidString(body.Message) {
			if dec, err := charmap.Windows1251.NewDecoder().String(body.Message); err == nil {
				body.Message = dec
			}
		}
		if body.Temp == 0 {
			body.Temp = defaultTemperature
		}
		lang := services.DetectLang(body.Message)
		w.Header().Set("X-Lang", lang)

		// история + system
		history := sessions.Get(id)
		if len(history) == 0 {
			systemPrompt := buildSystemPrompt(lang)
			history = append(history, services.ChatMessage{Role: "system", Content: systemPrompt})
		}

		// мини-RAG по продуктам
		prods := services.SearchProducts(body.Message, lang, 3)
		ctxBullets := services.BuildContextBullets(prods, lang)
		if len(ctxBullets) > 0 {
			history = append(history, services.ChatMessage{
				Role:    "assistant",
				Content: "Context:\n" + strings.Join(ctxBullets, "\n"),
			})
		}

		// факты из базы знаний
		know := services.SearchKnowledge(body.Message, 3)
		if len(know) > 0 {
			lines := []string{}
			for _, f := range know {
				lines = append(lines, "• "+f.Title+": "+f.Text)
			}
			history = append(history, services.ChatMessage{
				Role:    "system",
				Content: "Use these verified facts about Zaman Bank:\n" + strings.Join(lines, "\n"),
			})
		}

		// пользовательский запрос
		userMsg := services.ChatMessage{Role: "user", Content: body.Message}
		sessions.Append(id, userMsg)
		history = append(history, userMsg)
		history = trimHistoryIfNeeded(history)

		// LLM
		llm := services.NewLLM(cfgBaseURL, cfgAPIKey, cfgBackend)
		reply, err := llm.Chat(cfgModel, history, body.Temp)
		fallback := false
		if err != nil || strings.TrimSpace(reply) == "" {
			reply = services.BuildLocalAnswer(lang, body.Message, ctxBullets)
			fallback = true
		}

		// мягкая «рамка» под Zamanbank
		reply = services.RefineToZamanBank(reply, lang)
		services.LogFeedback(body.Message, reply)

		w.Header().Set("X-Fallback", map[bool]string{true: "true", false: "false"}[fallback])
		sessions.Append(id, services.ChatMessage{Role: "assistant", Content: reply})
		_ = json.NewEncoder(w).Encode(map[string]any{
			"session_id":  id,
			"reply":       reply,
			"history_len": len(sessions.Get(id)),
			"fallback":    fallback,
		})
	})

	// --- train & rag ---
	mux.HandleFunc("/admin/train", handlers.TrainHandler)
	mux.HandleFunc("/rag/chat", handlers.RagChatHandler)

	// --- статика ---
	mux.HandleFunc("/chat.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticRoot, "chat.html"))
	})
	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/chat.html", http.StatusFound)
	})
	mux.Handle("/", http.FileServer(http.Dir(staticRoot)))

	// middleware
	handler := withRequestID(withLogging(mux))
	log.Println("listening on :8080 (backend:", cfgBackend, ", model:", cfgModel, ")")
	log.Fatal(http.ListenAndServe(":8080", handler))
}