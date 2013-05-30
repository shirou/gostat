package outputs

import (
	"../record"
	"fmt"
	"strings"
)

type CSV struct{}

func (l CSV) Header(r record.Record) {}

func (l CSV) Emit(rs []record.Record, args []string) error {
	kv := make([]string, 0)
	keys := make([]string, 0)
	for _, r := range rs {
		for k, v := range r.Value {
			keys = append(keys, fmt.Sprintf("%s", k))
			kv = append(kv, fmt.Sprintf("%s", v))
		}
	}
	fmt.Println(strings.Join(keys, ","))
	fmt.Println(strings.Join(kv, ","))
	return nil
}
