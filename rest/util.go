package rest

import (
	"encoding/json"
)

// pretty pretty-prints an interface using the JSON marshaler
func pretty(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "\t")
	return string(b)
}
