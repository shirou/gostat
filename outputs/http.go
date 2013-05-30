package outputs

import (
	"net/url"
	"net/http"
	"bitbucket.org/r_rudi/gostat/record"
	"encoding/json"
)

type Http struct{}

func (l Http) Header(r record.Record) {}

func (l Http) Emit(rs []record.Record, args []string) error {
	dst := args[0]

	for _, r := range rs {
		value := make(map[string]string, 0)
		value["time"] = r.Time.UTC().String()
		value["tag"] = r.Tag

		for k, v := range r.Value {
			value[k] = v
		}

		data, err := json.Marshal(value)
		if err != nil{
			continue
		}

		values := url.Values{}
        values.Set("json", string(data))

		http.PostForm(dst, values)
	}
	return nil
}
