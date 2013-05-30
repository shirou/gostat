package outputs

import (
	"bitbucket.org/r_rudi/gostat/record"
)

type Output interface {
	Emit([]record.Record, []string) error
}
