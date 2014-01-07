package modules

import (
	"bitbucket.org/r_rudi/gostat/record"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Disk struct {
	Name  string
	Nick  []string
	Vars  []string
	Type  string
	Scale int
}

func NewDisk() Plugin {
	p := Disk{
		"disk",
		[]string{"read", "writ"},
		[]string{"read", "writ"},
		"d",
		34,
	}
	return Plugin(p)
}

func (p Disk) Check(conf map[string]map[string]string) error {
	filename := "/proc/diskstats"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	} else {
		return nil
	}
	return nil
}

func (p Disk) Extract(retchan chan record.Record, conf map[string]map[string]string) {
	fmt.Println(p.Name + " check")

	filename := "/proc/diskstats"
	s, err := ReadLines(filename)
	if err != io.EOF {
		return
	}

	vars := map[string]string{}
	values := strings.Fields(s[0]) // get only first line
	for i, v := range p.Vars {
		vars[v] = values[i+1]
	}

	r := record.Record{p.Name, time.Now(), vars}

	retchan <- r
	return
}
