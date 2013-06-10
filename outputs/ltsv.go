package outputs

import (
	"bitbucket.org/r_rudi/gostat/record"
	"fmt"
	"strings"
)

type LTSV struct{}

func (l LTSV) Header(r record.Record) {}

func (l LTSV) Emit(rs []record.Record, conf map[string]map[string]string) error {
	for _, r := range rs {
		buf := fmt.Sprintf("time:%s\ttag:%s\t", r.Time, r.Tag)
		kv := make([]string, 0)
		for k, v := range r.Value {
			kv = append(kv, fmt.Sprintf("%s:%s", k, v))
		}
		buf += strings.Join(kv, "\t")
		fmt.Println(buf)
	}

	return nil
}
