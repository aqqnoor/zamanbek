package services

import (
	"regexp"
	"strings"
)

// предкомпилируем несколько общих паттернов (короткий список, не «миллион исключений»)
var (
	reBanksGeneric = regexp.MustCompile(`(?i)\b(все\s+банки|банки\s+обычно|обычно\s+банки|most\s+banks|banks\s+usually|commercial\s+banks|коммерческ\w+\s+банк\w+|центральн\w+\s+банк)\b`)
	reCreditGeneric = regexp.MustCompile(`(?i)\b(кредит(ы|ов)?|loans?)\b`)
)

// RefineToZamanBank — мягко переписывает слишком общие ответы под контекст Zamanbank.
func RefineToZamanBank(reply, lang string) string {
	r := strings.TrimSpace(reply)
	if r == "" {
		return r
	}
	lower := strings.ToLower(r)

	// 1) убираем «все банки / банки обычно»
	if reBanksGeneric.FindStringIndex(lower) != nil {
		switch lang {
		case "kz":
			r = reBanksGeneric.ReplaceAllString(r, "Zamanbank-та біз")
		case "ru":
			r = reBanksGeneric.ReplaceAllString(r, "В Zamanbank мы")
		default:
			r = reBanksGeneric.ReplaceAllString(r, "At Zamanbank we")
		}
	}

	// 2) если речь про «кредиты» без упоминания халяль-эквивалентов — добавим пояснение
	if reCreditGeneric.FindStringIndex(lower) != nil &&
		!strings.Contains(lower, "мурабах") &&
		!strings.Contains(lower, "ijara") && !strings.Contains(lower, "ижар") &&
		!strings.Contains(lower, "wakala") && !strings.Contains(lower, "вакал") {
		switch lang {
		case "kz":
			r += "\n\nZamanbank дәстүрлі кредиттердің орнына халал қаржыландыру ұсынады: Мурабаха, Ижара, Вакала."
		case "ru":
			r += "\n\nВ Zamanbank вместо классических кредитов доступны халяль-формы финансирования: Мурабаха, Иджара, Вакала."
		default:
			r += "\n\nAt Zamanbank we use halal financing instead of conventional loans: Murabaha, Ijara, Wakala."
		}
	}

	// 3) подстрахуемся: если ответ вообще не содержит упоминания банка — мягко вставим «рамку»
	if !strings.Contains(strings.ToLower(r), "zaman") && !strings.Contains(strings.ToLower(r), "заман") {
		switch lang {
		case "kz":
			r = "Zamanbank контекстінде: " + r
		case "ru":
			r = "В контексте Zamanbank: " + r
		default:
			r = "In the context of Zamanbank: " + r
		}
	}

	return r
}
