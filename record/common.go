package record

import (
	"time"
)

// Resouce information. This is created one resouce by one.
type Record struct {
	Tag   string
	Time  time.Time
	Value map[string]string
}
