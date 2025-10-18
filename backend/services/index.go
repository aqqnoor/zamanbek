package services

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Doc struct {
	ID      string            `json:"id"`
	Title   string            `json:"title"`
	Content string            `json:"content"`
	Meta    map[string]string `json:"meta,omitempty"`
	Vector  []float32         `json:"vector"`
}

type KeywordBoost struct {
	Phrase   string   `json:"phrase"`
	Weight   float64  `json:"weight"`
	Synonyms []string `json:"synonyms"`
}
type BoostConfig struct {
	Boosts []KeywordBoost `json:"boosts"`
}

type FlatIndex struct {
	Docs   []Doc       `json:"docs"`
	Boosts BoostConfig `json:"boosts"`
}

func IndexPath() string { return filepath.Join("data", "index.json") }

func (fi *FlatIndex) Save(path string) error {
	b, _ := json.MarshalIndent(fi, "", "  ")
	return os.WriteFile(path, b, 0644)
}
func LoadIndex(path string) (*FlatIndex, error) {
	b, err := os.ReadFile(path)
	if err != nil { return &FlatIndex{}, nil }
	var fi FlatIndex
	if err := json.Unmarshal(b, &fi); err != nil { return nil, err }
	return &fi, nil
}

func cosine(a, b []float32) float64 {
	var dot, na, nb float64
	n := len(a); if len(b) < n { n = len(b) }
	for i := 0; i < n; i++ {
		da := float64(a[i]); db := float64(b[i])
		dot += da*db; na += da*da; nb += db*db
	}
	if na == 0 || nb == 0 { return 0 }
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func (fi *FlatIndex) boost(query, text string) float64 {
	q := strings.ToLower(query)
	d := strings.ToLower(text)
	var s float64
	for _, b := range fi.Boosts.Boosts {
		hit := strings.Contains(q, strings.ToLower(b.Phrase)) || strings.Contains(d, strings.ToLower(b.Phrase))
		if !hit {
			for _, syn := range b.Synonyms {
				if strings.Contains(q, strings.ToLower(syn)) || strings.Contains(d, strings.ToLower(syn)) {
					hit = true; break
				}
			}
		}
		if hit { s += b.Weight }
	}
	return s
}

type ScoredDoc struct {
	Doc   Doc     `json:"doc"`
	Score float64 `json:"score"`
}

func (fi *FlatIndex) Search(vec []float32, query string, k int) []ScoredDoc {
	best := make([]ScoredDoc, 0, k)
	for _, d := range fi.Docs {
		base := cosine(d.Vector, vec)
		boost := fi.boost(query, d.Title+" "+d.Content)
		score := base + boost

		ins := false
		for i := range best {
			if score > best[i].Score {
				best = append(best[:i+1], best[i:]...)
				best[i] = ScoredDoc{Doc: d, Score: score}
				ins = true
				break
			}
		}
		if !ins && len(best) < k {
			best = append(best, ScoredDoc{Doc: d, Score: score})
		}
		if len(best) > k { best = best[:k] }
	}
	return best
}