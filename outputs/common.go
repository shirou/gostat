package outputs

import (
	"bitbucket.org/r_rudi/gostat/record"
)

// the interface of Output modules
type Output interface {
	Emit([]record.Record, map[string]map[string]string) error
}
