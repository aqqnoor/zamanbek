package services

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

// normalizeLabel приводит заголовки к единым ключам.
func normalizeLabel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", " ") // NBSP → space
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")

	switch s {
	case "продукт 1", "product 1", "өнім 1":
		return "name"

	case "тип продукта", "product type", "өнім түрі":
		return "type"

	case "наценка (тг)", "размер наценки, тенге", "markup", "markup (kzt)":
		return "markup"

	case "максимальная сумма, тенге", "максимальная сумма", "max amount", "max amount (kzt)":
		return "max_sum"

	case "минимальная сумма, тенге", "минимальная сумма", "min amount", "min amount (kzt)":
		return "min_sum"

	case "описание", "description":
		return "description"

	default:
		return s // неизвестные ключи не теряем — положим в Extra
	}
}

// ExcelProduct — промежуточная структура результата парсинга одного продукта из Excel.
type ExcelProduct struct {
	Sheet       string            // Название листа (Retail/SME/…)
	RawFields   map[string]string // Нормализованные поля (name/type/markup/min_sum/max_sum/description/…)
	RawOriginal map[string]string // Исходные лейблы → значения (на всякий случай)
}

// ParseExcelProducts читает Excel и возвращает список продуктов, где каждые несколько строк — один продукт.
func ParseExcelProducts(path string) ([]ExcelProduct, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open excel: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheets := f.GetSheetList()
	out := make([]ExcelProduct, 0, 32)

	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			log.Printf("[excel] skip sheet %s: %v", sheet, err)
			continue
		}

		var cur *ExcelProduct

		for _, row := range rows {
			// ожидаем, что лейбл в колонке B (index 1), значение — в колонке C (index 2)
			var label, value string
			if len(row) > 1 {
				label = strings.TrimSpace(row[1])
			}
			if len(row) > 2 {
				value = strings.TrimSpace(row[2])
			}
			if label == "" && value == "" {
				continue
			}

			norm := normalizeLabel(label)
			if norm == "name" {
				// начинаем новый продукт
				if cur != nil && len(cur.RawFields) > 0 {
					out = append(out, *cur)
				}
				cur = &ExcelProduct{
					Sheet:       sheet,
					RawFields:   map[string]string{},
					RawOriginal: map[string]string{},
				}
				if value != "" {
					cur.RawFields["name"] = value
				}
				cur.RawOriginal[label] = value
				continue
			}

			// если продукта ещё нет, но попались строки — создаём контейнер
			if cur == nil {
				cur = &ExcelProduct{
					Sheet:       sheet,
					RawFields:   map[string]string{},
					RawOriginal: map[string]string{},
				}
			}

			if norm != "" && value != "" {
				cur.RawFields[norm] = value
			}
			if label != "" {
				cur.RawOriginal[label] = value
			}
		}
		// финальный хвост
		if cur != nil && len(cur.RawFields) > 0 {
			out = append(out, *cur)
		}
	}

	return out, nil
}

// canonID строит предсказуемый ID на основе имени.
func canonID(name string) string {
	id := strings.ToLower(name)
	id = strings.TrimSpace(id)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "_", "-")
	id = regexp.MustCompile(`[^a-z0-9\-]+`).ReplaceAllString(id, "")
	if id == "" {
		id = "product"
	}
	return id
}

// ExcelToProducts конвертирует ExcelProduct → Product (наша унифицированная модель).
func ExcelToProducts(xs []ExcelProduct) []Product {
	out := make([]Product, 0, len(xs))
	for _, x := range xs {
		name := strings.TrimSpace(x.RawFields["name"])
		typ := strings.TrimSpace(x.RawFields["type"])
		markup := strings.TrimSpace(x.RawFields["markup"])
		minSum := strings.TrimSpace(x.RawFields["min_sum"])
		maxSum := strings.TrimSpace(x.RawFields["max_sum"])
		desc := strings.TrimSpace(x.RawFields["description"])

		tags := []string{}
		if name != "" {
			tags = append(tags, name)
		}
		if typ != "" {
			tags = append(tags, typ)
		}
		if markup != "" {
			tags = append(tags, markup)
		}

		p := Product{
			ID:        canonID(name),
			Category:  x.Sheet,     // Retail/SME/…
			Name:      name,
			Type:      typ,
			Markup:    markup,
			MinSum:    minSum,
			MaxSum:    maxSum,
			Halal:     true,        // по умолчанию считаем халяль (можно добавить явное поле в Excel)
			Tags:      tags,
			Extra:     x.RawFields, // чтобы ничего не потерять
			RawSource: filepath.Base(x.Sheet),
			Desc:      desc,
		}
		out = append(out, p)
	}
	return out
}
