package services

import "testing"

func TestDetectLang(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Сәлем! Бұл карта халал ма?", "kz"},
		{"Привет! Эта карта халяльная?", "ru"},
		{"Hi! Is this card halal?", "en"},
		{"", "en"}, // пустое → en по умолчанию
		{"12345", "en"},
		{"Қазақша мен Русский аралас", "kz"},
		{"Русский и English mixed", "ru"},
	}
	for _, tt := range tests {
		got := DetectLang(tt.in)
		if got != tt.want {
			t.Errorf("DetectLang(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
