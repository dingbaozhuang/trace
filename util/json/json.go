// +build !jsoniter

package json

import "encoding/json"

// stand library json
var (
	Marshal   = json.Marshal
	Unmarshal = json.Unmarshal
)
