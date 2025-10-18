package handlers

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "path/filepath"
)

// TrainHandler loads the religion keyword boosts from the JSON file on disk
// and returns summary information about the number of phrases and total
// synonyms.
func TrainHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    // The file resides in data/religion_keywords.json when invoked from backend/.
    path := filepath.Join("data", "religion_keywords.json")
    data, err := ioutil.ReadFile(path)
    if err != nil {
        http.Error(w, fmt.Sprintf("failed to read keywords file: %v", err), http.StatusInternalServerError)
        return
    }
    var payload struct {
        Boosts []struct {
            Phrase   string   `json:"phrase"`
            Weight   float64  `json:"weight"`
            Synonyms []string `json:"synonyms"`
        } `json:"boosts"`
    }
    if err := json.Unmarshal(data, &payload); err != nil {
        http.Error(w, fmt.Sprintf("failed to parse keywords JSON: %v", err), http.StatusInternalServerError)
        return
    }
    totalSynonyms := 0
    for _, boost := range payload.Boosts {
        totalSynonyms += len(boost.Synonyms)
    }
    resp := map[string]any{
        "phrases":       len(payload.Boosts),
        "total_synonyms": totalSynonyms,
        "status":        "loaded",
    }
    _ = json.NewEncoder(w).Encode(resp)
}