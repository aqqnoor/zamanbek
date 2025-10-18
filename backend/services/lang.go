package services

import (
	"unicode"
)

// DetectLang определяет язык сообщения.
// 1) Если есть хотя бы один казахский символ — "kz"
// 2) Иначе если есть любой кириллический символ — "ru"
// 3) Иначе — "en"
func DetectLang(s string) string {
	// Казахский набор (доп. к кириллице)
	kazakhRunes := map[rune]bool{
		'Ә': true, 'ә': true, 'Ғ': true, 'ғ': true, 'Қ': true, 'қ': true,
		'Ң': true, 'ң': true, 'Ө': true, 'ө': true, 'Ұ': true, 'ұ': true,
		'Ү': true, 'ү': true, 'Һ': true, 'һ': true, 'І': true, 'і': true,
	}
	hasCyr := false
	for _, r := range s {
		if kazakhRunes[r] {
			return "kz"
		}
		if unicode.In(r, unicode.Cyrillic) {
			hasCyr = true
		}
	}
	if hasCyr {
		return "ru"
	}
	return "en"
}