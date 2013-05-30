package record

import (
	"time"
)

type Record struct {
	Tag   string
	Time  time.Time
	Value map[string]string
}
