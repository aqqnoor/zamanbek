package services

import "strings"

// IsUnhelpful определяет, бесполезен ли ответ модели (отмазка/болтовня).
func IsUnhelpful(reply, lang string) bool {
	r := strings.ToLower(strings.TrimSpace(reply))
	if r == "" {
		return true
	}
	// короткие "не понимаю" во всех трёх языках
	badPhrases := []string{
		"я не понимаю", "не понимаю ваш запрос", "уточните", "не могу понять",
		"i don't understand", "i do not understand", "please clarify",
		"түсінбедім", "түсінбеймін", "нақтылаңыз", "анығырақ айтыңыз",
	}
	for _, p := range badPhrases {
		if strings.Contains(r, p) && len([]rune(r)) < 200 {
			return true
		}
	}
	// Одно-два предложения без фактики при наличии контекста тоже считаем слабым ответом:
	if len([]rune(r)) < 40 {
		return true
	}
	return false
}