package services

import (
    "fmt"
    "strings"
)

// LLMClient — общее интерфейсное API обоих клиентов.
type LLMClient interface {
    Chat(model string, messages []ChatMessage, temperature float32) (string, error)
}

// MockLLM — локальная заглушка для офлайн-режима.
type MockLLM struct{}

func NewLLMMock() *MockLLM { return &MockLLM{} }

func (l *MockLLM) Chat(model string, messages []ChatMessage, temperature float32) (string, error) {
    if len(messages) == 0 {
        return "", fmt.Errorf("no messages provided")
    }
    userMsg := strings.TrimSpace(messages[len(messages)-1].Content)
    lang := DetectLang(userMsg)
    if userMsg == "" {
        switch lang {
        case "kz":
            return "Сіз сұрақ енгізбедіңіз. Қалай көмектесе аламын?", nil
        case "ru":
            return "Вы не ввели вопрос. Как я могу помочь?", nil
        default:
            return "You didn't type a question. How can I help?", nil
        }
    }
    switch lang {
    case "kz":
        return fmt.Sprintf("Құрметті клиент, сұрағыңызды қабылдадық: %s", userMsg), nil
    case "ru":
        return fmt.Sprintf("Уважаемый клиент, ваш вопрос принят: %s", userMsg), nil
    default:
        return fmt.Sprintf("Dear customer, we received your question: %s", userMsg), nil
    }
}