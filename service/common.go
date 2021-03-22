package service

import (
	"bytes"
	"encoding/json"
)

//JsonFormat Json outupt.
func JsonFormat(v interface{}) string {
	// if msg == "" {
	// 	return ""
	// }
	var out bytes.Buffer

	bs, _ := json.Marshal(v)
	json.Indent(&out, bs, "", "\t")

	return out.String()
}
