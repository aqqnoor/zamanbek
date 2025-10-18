package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// BankLLM — HTTP-клиент к роутеру банка (OpenAI-совместимый /chat/completions).
type BankLLM struct {
    baseURL string
    apiKey  string
    httpc   *http.Client
}

func NewLLMBank(baseURL, apiKey string) *BankLLM {
    return &BankLLM{
        baseURL: trimTrailingSlash(baseURL),
        apiKey:  apiKey,
        httpc: &http.Client{
            Timeout: 15 * time.Second,
        },
    }
}

func trimTrailingSlash(s string) string {
    if len(s) > 0 && s[len(s)-1] == '/' {
        return s[:len(s)-1]
    }
    return s
}

type oaChatReq struct {
    Model       string        `json:"model"`
    Temperature float32       `json:"temperature,omitempty"`
    Messages    []ChatMessage `json:"messages"`
}

type oaChoice struct {
    Message ChatMessage `json:"message"`
}

type oaChatResp struct {
    Choices []oaChoice `json:"choices"`
}

func (l *BankLLM) Chat(model string, messages []ChatMessage, temperature float32) (string, error) {
    // Формируем OpenAI-совместимый запрос.
    payload := oaChatReq{
        Model:       model,
        Temperature: temperature,
        Messages:    messages,
    }
    body, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", l.baseURL+"/v1/chat/completions", bytes.NewReader(body))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")
    if l.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+l.apiKey)
    }
    req.Header.Set("X-Client", "zamanbek-assistant")

    // Мини-ретраи.
    var lastErr error
    for i := 0; i < 2; i++ {
        resp, err := l.httpc.Do(req)
        if err != nil {
            lastErr = err
            continue
        }
        defer resp.Body.Close()
        if resp.StatusCode < 200 || resp.StatusCode >= 300 {
            lastErr = fmt.Errorf("bank LLM http %d", resp.StatusCode)
            continue
        }
        var out oaChatResp
        if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
            lastErr = err
            continue
        }
        if len(out.Choices) == 0 {
            return "", fmt.Errorf("empty choices")
        }
        return out.Choices[0].Message.Content, nil
    }
    return "", lastErr
}
