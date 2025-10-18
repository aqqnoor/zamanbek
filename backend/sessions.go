package main

import (
    "crypto/rand"
    "encoding/hex"
    "net/http"
    "sync"

    "zmb-assistant/services"
)

type SessionStore struct {
    mu       sync.RWMutex
    sessions map[string][]services.ChatMessage
}

func NewSessionStore() *SessionStore {
    return &SessionStore{
        sessions: make(map[string][]services.ChatMessage),
    }
}

func (s *SessionStore) ensureID(r *http.Request) (string, bool) {
    if cookie, err := r.Cookie("session_id"); err == nil && cookie.Value != "" {
        return cookie.Value, true
    }
    return generateSessionID(), false
}

func (s *SessionStore) EnsureSessionID(w http.ResponseWriter, r *http.Request) string {
    id, exists := s.ensureID(r)
    if !exists {
        http.SetCookie(w, &http.Cookie{Name: "session_id", Value: id, Path: "/"})
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    if _, ok := s.sessions[id]; !ok {
        s.sessions[id] = []services.ChatMessage{}
    }
    return id
}

func (s *SessionStore) Reset(id string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.sessions[id] = []services.ChatMessage{}
}

func (s *SessionStore) Get(id string) []services.ChatMessage {
    s.mu.RLock()
    defer s.mu.RUnlock()
    hist, ok := s.sessions[id]
    if !ok {
        return nil
    }
    cp := make([]services.ChatMessage, len(hist))
    copy(cp, hist)
    return cp
}

func (s *SessionStore) Append(id string, msg services.ChatMessage) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.sessions[id] = append(s.sessions[id], msg)
}

func generateSessionID() string {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        return ""
    }
    return hex.EncodeToString(b)
}