package modules

import (
	"bitbucket.org/r_rudi/gostat/record"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Cpu struct {
	Name  string
	Nick  []string
	Vars  []string
	Type  string
	Scale int
}

func NewCpu() (Plugin, error) {
	p := Cpu{
		"cpu",
		[]string{"usr", "sys", "idl", "wai", "hiq", "siq", "stl"}, // 2.6.11 or later
		[]string{"usr", "sys", "idl", "wai", "hiq", "siq", "stl"}, // 2.6.11 or later
		"p",
		34,
	}
	return Plugin(p), p.Check()
}

func (p Cpu) Check() error {
	filename := "/proc/stat"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	} else {
		return nil
	}
	return nil  // for 1.0
}

func (p Cpu) GetVal(filename string, vars []string) ([]int, float32, error) {
	s, err := ReadLines(filename)
	if err != io.EOF {
		return []int{}, 0, err
	}

	values := strings.Fields(s[0]) // only first line

	ret := []int{}
	total := 0
	for i, _ := range vars {
		i, err := strconv.Atoi(values[i+1])
		if err == nil {
			total = total + i
			ret = append(ret, i)
		}
	}

	return ret, float32(total), nil
}

func (p Cpu) Extract(retchan chan record.Record) {
	filename := "/proc/stat"
	s1, total1, err := p.GetVal(filename, p.Vars)
	if err != nil {
		close(retchan)
		return
	}
	time.Sleep(time.Millisecond * 200) // FIXME: why 100?
	s2, total2, err := p.GetVal(filename, p.Vars)
	if err != nil {
		close(retchan)
		return
	}

	vars := map[string]string{}
	for i, _ := range s2 {
		title := p.Vars[i]
		//		fmt.Println(s2[i],s1[i], s2[i]-s1[i], (total2 - total1))
		var result float32
		result = (float32(s2[i]-s1[i]) / (total2 - total1)) * 100
		vars[title] = fmt.Sprintf("%0.2f", result)
	}

	r := record.Record{p.Name, time.Now(), vars}

	retchan <- r
	return
}
