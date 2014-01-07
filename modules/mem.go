package modules

import (
	"bitbucket.org/r_rudi/gostat/record"
	"io"
	"os"
	"strings"
	"time"
)

type Mem struct {
	Name  string
	Nick  []string
	Vars  []string
	Type  string
	Scale int
}

func NewMem() Plugin {
	p := Mem{
		"memory_usage",
		[]string{"used", "buff", "cach", "free"},
		[]string{"MemUsed", "Buffers", "Cached", "MemFree"},
		"f",
		1024,
	}
	return Plugin(p)
}

func (p Mem) Check(conf map[string]map[string]string) error {
	filename := "/proc/meminfo"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	} else {
		return nil
	}
	return nil
}

func (p Mem) Extract(retchan chan record.Record, conf map[string]map[string]string) {
	filename := "/proc/meminfo"
	s, err := ReadLines(filename)
	if err != io.EOF {
		return
	}

	ret := map[string]string{}
	for _, line := range s {
		values := strings.Fields(line)
		for _, val := range p.Vars {
			if strings.Replace(values[0], ":", "", 1) == val {
				ret[val] = values[1]
			}
		}
	}

	r := record.Record{p.Name, time.Now(), ret}
	retchan <- r
	return
}
