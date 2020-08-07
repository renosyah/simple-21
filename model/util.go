package model

import "encoding/json"

func ToJson(e interface{}) []byte {
	jbyte, _ := json.Marshal(e)
	return jbyte
}
