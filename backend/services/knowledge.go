package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Fact struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

var knowledgeBase []Fact

func LoadKnowledge() {
	path := "data/knowledge.json"
	if _, err := os.Stat(path); err != nil {
		log.Printf("[knowledge] no file: %v", err)
		return
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("[knowledge] read err: %v", err)
		return
	}
	var wrap struct {
		Facts []Fact `json:"facts"`
	}
	if err := json.Unmarshal(data, &wrap); err != nil {
		log.Printf("[knowledge] parse err: %v", err)
		return
	}
	knowledgeBase = wrap.Facts
	log.Printf("[knowledge] loaded %d facts", len(knowledgeBase))
}

// простой поиск по совпадениям слов
func SearchKnowledge(q string, limit int) []Fact {
	if limit <= 0 {
		limit = 3
	}
	qlow := strings.ToLower(q)
	type scored struct {
		f Fact
		s float64
	}
	var out []scored
	for _, f := range knowledgeBase {
		text := strings.ToLower(f.Title + " " + f.Text)
		words := strings.Fields(qlow)
		score := 0.0
		for _, w := range words {
			if strings.Contains(text, w) {
				score++
			}
		}
		if len(words) > 0 {
			score = score / float64(len(words))
		}
		if score > 0 {
			out = append(out, scored{f, score})
		}
	}
	// простая сортировка пузырьком
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].s > out[i].s {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	res := []Fact{}
	for i, s := range out {
		if i >= limit {
			break
		}
		res = append(res, s.f)
	}
	return res
}

func AddFactFromFeedback(entry FeedbackEntry) {
	if strings.TrimSpace(entry.Query) == "" || strings.TrimSpace(entry.Reply) == "" {
		return
	}
	knowledgeBase = append(knowledgeBase, Fact{
		ID:    "fb_" + strings.ReplaceAll(entry.Query, " ", "_"),
		Title: entry.Query,
		Text:  entry.Reply,
	})
	log.Printf("[selftrain] added fact from feedback: %s", entry.Query)
}
