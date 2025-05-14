package utils

import (
	"fmt"
	"sort"
	"strings"
)

// Message - contain key / value for response
type Message map[string]any

// create line to response format {key1:value2},{key2,value2}...
// sort.Strings make the string predictable
func (msg Message) String() string {
	lineMsg := make([]string, 0, len(msg))
	for k, v := range msg {
		lineMsg = append(lineMsg, fmt.Sprintf(`{%s:%v}`, k, v))
	}
	sort.Strings(lineMsg)
	return strings.Join(lineMsg, ",")
}
