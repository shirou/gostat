package outputs

import (
	"bitbucket.org/r_rudi/gostat/record"
	"fmt"
	"strings"
)

type WhiteSpace struct{}

func (l WhiteSpace) Header(r record.Record) {}

func (l WhiteSpace) Emit(rs []record.Record, conf map[string]map[string]string) error {
	kv := make([]string, 0)
	keys := make([]string, 0)
	for _, r := range rs {
		for k, v := range r.Value {
			keys = append(keys, fmt.Sprintf("%s", k))
			kv = append(kv, fmt.Sprintf("%s", v))
		}
	}
	fmt.Println(strings.Join(keys, " "))
	fmt.Println(strings.Join(kv, " "))
	return nil
}
