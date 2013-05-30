package outputs

import (
	"../record"
)

type Output interface {
	Emit([]record.Record, []string) error
}
