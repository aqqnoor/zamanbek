package jobs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"zmb-assistant/services"
)

type RawDoc struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func readJSON[T any](path string, out *T) error {
	b, err := os.ReadFile(path)
	if err != nil { return err }
	return json.Unmarshal(b, out)
}

func chunk(id, title, content string) []services.Doc {
	const sz = 3500 // грубый, но быстрый чанк по символам
	var docs []services.Doc
	for i := 0; i < len(content); i += sz {
		j := i + sz; if j > len(content) { j = len(content) }
		docs = append(docs, services.Doc{
			ID: fmt.Sprintf("%s_%d", id, i/sz),
			Title: title,
			Content: content[i:j],
			Meta: map[string]string{"src": id},
		})
	}
	return docs
}

func TrainIndex(baseURL, apiKey string) error {
	// 1) читаем все источники
	var products, goals, faqs []RawDoc
	_ = readJSON(filepath.Join("data","products.json"), &products)
	_ = readJSON(filepath.Join("data","goals.json"),    &goals)
	_ = readJSON(filepath.Join("data","faq.json"),      &faqs)

	all := make([]services.Doc, 0, len(products)+len(goals)+len(faqs))
	for _, r := range [][]RawDoc{products, goals, faqs} {
		for _, d := range r {
			all = append(all, chunk(d.ID, d.Title, d.Content)...)
		}
	}
	if len(all) == 0 {
		return fmt.Errorf("no docs to index (products/goals/faq are empty?)")
	}

	// 2) батч-эмбеддинги
	ec := services.NewEmbedding(baseURL, apiKey)
	texts := make([]string, len(all))
	for i := range all {
		texts[i] = all[i].Title + "\n" + all[i].Content
	}
	vecs, err := ec.EmbedMany(texts)
	if err != nil { return err }
	for i := range all { all[i].Vector = vecs[i] }

	// 3) бусты религиозных терминов
	var boosts services.BoostConfig
	_ = readJSON(filepath.Join("data","religion_keywords.json"), &boosts)

	// 4) сохраняем индекс
	idx := &services.FlatIndex{Docs: all, Boosts: boosts}
	return idx.Save(services.IndexPath())
}