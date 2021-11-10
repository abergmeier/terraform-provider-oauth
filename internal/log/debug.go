package log

import (
	"encoding/json"
	"log"
	"sort"
)

func DebugLogJSON(j []byte) {
	m := make(map[string]interface{})
	err := json.Unmarshal(j, &m)
	if err != nil {
		panic(err)
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	log.Printf(`[DEBUG] Found keys in auth file: %s\n
`, keys)
}
