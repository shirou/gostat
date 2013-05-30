package modules

import (
	"bitbucket.org/r_rudi/gostat/record"
	"io"
	"os"
	"strings"
	"time"
)

type Load struct {
	Name  string
	Nick  []string
	Vars  []string
	Type  string
	Scale float32
}

func NewLoad() (Plugin, error) {
	p := Load{
		"load",
		[]string{"1m", "5m", "15m"},
		[]string{"load1", "load5", "load15"},
		"f",
		0.5,
	}
	return Plugin(p), p.Check()
}

func (p Load) Check() error {
	filename := "/proc/loadavg"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	} else {
		return nil
	}
	return nil // for 1.0
}

func (p Load) Extract(retchan chan record.Record) {
	filename := "/proc/loadavg"
	s, err := ReadLines(filename)
	if err != io.EOF {
		close(retchan)
		return
	}

	values := strings.Fields(s[0])
	ret := map[string]string{}
	for i, t := range p.Vars {
		ret[t] = values[i]
	}

	r := record.Record{p.Name, time.Now(), ret}
	retchan <- r
	return
}
