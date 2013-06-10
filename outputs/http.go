package outputs

import (
	"bitbucket.org/r_rudi/gostat/record"
	"encoding/json"
	"net/http"
	"net/url"
)

type Http struct{}

func (l Http) Header(r record.Record) {}

func (l Http) Emit(rs []record.Record, conf map[string]map[string]string) error {
	dst := conf["root"]["server"]

	for _, r := range rs {
		value := make(map[string]string, 0)
		value["time"] = r.Time.UTC().String()
		value["tag"] = r.Tag

		for k, v := range r.Value {
			value[k] = v
		}

		data, err := json.Marshal(value)
		if err != nil {
			continue
		}

		values := url.Values{}
		values.Set("json", string(data))

		http.PostForm(dst, values)
	}
	return nil
}
