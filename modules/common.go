package modules

import (
	"../record"
	"bufio"
	"os"
	"strings"
)

type Plugin interface {
	Check() error
	Extract(chan record.Record)
}

func ReadLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	ret := make([]string, 0)

	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	for err == nil {
		ret = append(ret, strings.Trim(line, "\n"))
		line, err = r.ReadString('\n')
	}

	return ret, err
}
