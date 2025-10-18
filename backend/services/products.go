package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Product — унифицированная модель продукта для мини-RAG.
type Product struct {
	ID        string            `json:"id"`
	Category  string            `json:"category"` // Retail/SME/…
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Markup    string            `json:"markup,omitempty"`
	MinSum    string            `json:"min_sum,omitempty"`
	MaxSum    string            `json:"max_sum,omitempty"`
	Halal     bool              `json:"halal"`
	Tags      []string          `json:"tags,omitempty"`
	Desc      string            `json:"description,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
	RawSource string            `json:"raw_source,omitempty"` // sheet name
}

type ProductDB struct {
	Products []Product `json:"products"`
}

var prodDB ProductDB

// CountProducts — для /readyz
func CountProducts() int {
	return len(prodDB.Products)
}

// LoadProducts пытается загрузить продукты из Excel, иначе — из JSON.
func LoadProducts() {
	excelPath := "/mnt/data/Справочник по продуктам Zamanbank.xlsx"
	if st, err := os.Stat(excelPath); err == nil && !st.IsDir() {
		log.Printf("[products] loading from excel: %s", excelPath)
		if xs, err := ParseExcelProducts(excelPath); err == nil {
			prods := ExcelToProducts(xs)
			prodDB = ProductDB{Products: prods}
			log.Printf("[products] excel parsed: %d products", len(prods))
			return
		} else {
			log.Printf("[products] excel parse error: %v", err)
		}
	}

	// Фоллбэк: data/products.json
	jsonPath := filepath.Join("data", "products.json")
	if st, err := os.Stat(jsonPath); err == nil && !st.IsDir() {
		log.Printf("[products] loading from json: %s", jsonPath)
		data, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Printf("[products] read json error: %v", err)
			return
		}
		// Поддержка старого формата
		var tmp struct {
			Products []struct {
				ID       string              `json:"id"`
				Lang     []string            `json:"lang"`
				Type     string              `json:"type"`
				Names    map[string]string   `json:"names"`
				Halal    bool                `json:"halal"`
				Bullets  map[string][]string `json:"bullets"`
				Keywords []string            `json:"keywords"`
			} `json:"products"`
		}
		if err := json.Unmarshal(data, &tmp); err == nil && len(tmp.Products) > 0 {
			out := make([]Product, 0, len(tmp.Products))
			for _, p := range tmp.Products {
				name := p.Names["ru"]
				if name == "" {
					for _, v := range p.Names {
						name = v
						if name != "" {
							break
						}
					}
				}
				out = append(out, Product{
					ID:       p.ID,
					Category: "json",
					Name:     name,
					Type:     p.Type,
					Halal:    p.Halal,
					Tags:     append([]string{}, p.Keywords...),
					Extra:    map[string]string{},
					Desc:     "",
				})
			}
			prodDB = ProductDB{Products: out}
			log.Printf("[products] json parsed (legacy): %d products", len(out))
			return
		}
		var db ProductDB
		if err := json.Unmarshal(data, &db); err == nil && len(db.Products) > 0 {
			prodDB = db
			log.Printf("[products] json parsed: %d products", len(db.Products))
			return
		}
		log.Printf("[products] json parse error: file format not recognized")
		return
	}

	log.Printf("[products] no excel or json products found (excel: %s, json: %s)", excelPath, jsonPath)
}

// SearchProducts — простой поиск по имени, типу и тэгам.
func SearchProducts(q string, lang string, limit int) []Product {
	if limit <= 0 {
		limit = 3
	}
	qlow := strings.ToLower(q)
	out := make([]Product, 0, limit)

	for _, p := range prodDB.Products {
		hit := false
		if strings.Contains(strings.ToLower(p.Name), qlow) {
			hit = true
		}
		if !hit && strings.Contains(strings.ToLower(p.Type), qlow) {
			hit = true
		}
		if !hit {
			for _, kw := range p.Tags {
				if kw != "" && strings.Contains(strings.ToLower(kw), qlow) {
					hit = true
					break
				}
			}
		}
		if hit {
			out = append(out, p)
			if len(out) >= limit {
				break
			}
		}
	}
	return out
}

// BuildContextBullets — формирует короткий факт-контекст на нужном языке.
func BuildContextBullets(ps []Product, lang string) []string {
	bullets := []string{}
	for _, p := range ps {
		name := p.Name
		if name != "" {
			bullets = append(bullets, "• "+name)
		}
		if p.Type != "" {
			switch lang {
			case "kz":
				bullets = append(bullets, "  - Түрі: "+p.Type)
			case "ru":
				bullets = append(bullets, "  - Тип: "+p.Type)
			default:
				bullets = append(bullets, "  - Type: "+p.Type)
			}
		}
		if p.Markup != "" {
			switch lang {
			case "kz":
				bullets = append(bullets, "  - Үстеме: "+p.Markup)
			case "ru":
				bullets = append(bullets, "  - Наценка: "+p.Markup)
			default:
				bullets = append(bullets, "  - Markup: "+p.Markup)
			}
		}
		if p.MinSum != "" {
			switch lang {
			case "kz":
				bullets = append(bullets, "  - Минималды сома: "+p.MinSum)
			case "ru":
				bullets = append(bullets, "  - Минимальная сумма: "+p.MinSum)
			default:
				bullets = append(bullets, "  - Min amount: "+p.MinSum)
			}
		}
		if p.MaxSum != "" {
			switch lang {
			case "kz":
				bullets = append(bullets, "  - Максималды сома: "+p.MaxSum)
			case "ru":
				bullets = append(bullets, "  - Максимальная сумма: "+p.MaxSum)
			default:
				bullets = append(bullets, "  - Max amount: "+p.MaxSum)
			}
		}
		if p.Halal {
			switch lang {
			case "kz":
				bullets = append(bullets, "  - Халал талаптарына сай")
			case "ru":
				bullets = append(bullets, "  - Соответствует принципам халяль")
			default:
				bullets = append(bullets, "  - Sharia-compliant (halal)")
			}
		}
		if p.Desc != "" {
			switch lang {
			case "kz":
				bullets = append(bullets, "  - Сипаттама: "+p.Desc)
			case "ru":
				bullets = append(bullets, "  - Описание: "+p.Desc)
			default:
				bullets = append(bullets, "  - Description: "+p.Desc)
			}
		}
	}
	return bullets
}
