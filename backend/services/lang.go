// services/lang.go
package services

import "strings"

func IsKazakh(s string) bool {
    s = strings.ToLower(s)
    kzLetters := []string{"ә","ө","ү","ұ","қ","ғ","ң","і","һ"}
    for _, ch := range kzLetters {
        if strings.Contains(s, ch) { return true }
    }
    // эвристика по словам
    kzWords := []string{"ма","ме","ба","бе","па","пе","қалай","қандай","мүмкін","түсіндір"}
    for _, w := range kzWords {
        if strings.Contains(s, " "+w) { return true }
    }
    return false
}
