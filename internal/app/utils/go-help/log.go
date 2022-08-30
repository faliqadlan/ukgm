package help

import (
	"fmt"
	"time"
)

// Log func
func Log(message string) string {
	t := time.Now()
	return (fmt.Sprintf("%v %v", t.Format("["+time.RFC1123+"]"), message))
}
