// services/nlu.go
package services

import "strings"

type Intent string
const (
    IntentHalalCheck Intent = "halal_check"
    IntentGeneric    Intent = "generic"
)

func DetectIntent(q string) Intent {
    s := strings.ToLower(q)
    // любые формы халал/haram/риба → трактуем как проверку на халал
    keys := []string{"халал", "halal", "халяль", "haram", "харам", "риба", "өсім"}
    for _, k := range keys {
        if strings.Contains(s, k) {
            return IntentHalalCheck
        }
    }
    return IntentGeneric
}
