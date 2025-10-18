package services

import "strings"

// BuildLocalAnswer — гарантированный локальный ответ на языке пользователя.
// Если есть контекст по продуктам — кратко подставим bullets.
func BuildLocalAnswer(lang string, userMsg string, contextBullets []string) string {
    msg := strings.TrimSpace(userMsg)
    if lang == "kz" {
        if msg == "" {
            return "Сіз сұрақ енгізбедіңіз. Қалай көмектесе аламын?"
        }
        if len(contextBullets) > 0 {
            return "Сұрағыңыз қабылданды. Төменде пайдалы деректер:\n" + strings.Join(contextBullets, "\n")
        }
        return "Сұрағыңыз қабылданды. Қалай нақтылай аламын?"
    }
    if lang == "ru" {
        if msg == "" {
            return "Вы не ввели вопрос. Как я могу помочь?"
        }
        if len(contextBullets) > 0 {
            return "Ваш вопрос принят. Ниже полезная информация:\n" + strings.Join(contextBullets, "\n")
        }
        return "Ваш вопрос принят. Чем могу уточнить?"
    }
    // en
    if msg == "" {
        return "You didn't type a question. How can I help?"
    }
    if len(contextBullets) > 0 {
        return "Got it. Here are some helpful details:\n" + strings.Join(contextBullets, "\n")
    }
    return "Got it. How can I clarify further?"
}
