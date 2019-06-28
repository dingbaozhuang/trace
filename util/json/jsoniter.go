// +build jsoniter

package json

import "github.com/json-iterator/go"

// jsoniter
var (
	json = jsoniter.ConfigFastest

	Marshal   = json.Marshal
	Unmarshal = json.Unmarshal
)
