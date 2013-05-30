package modules

import (
	"bitbucket.org/r_rudi/gostat/record"
	"io"
	"os"
	"time"
)

type Aio struct {
	Name  string
	Nick  []string
	Vars  []string
	Type  string
	Scale int
}

func NewAio() (Plugin, error) {
	p := Aio{
		"aio",
		[]string{"#aio"},
		[]string{"aio"},
		"d",
		1024,
	}
	return Plugin(p), p.Check()
}

func (p Aio) Check() error {
	filename := "/proc/sys/fs/aio-nr"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return err
	} else {
		return nil
	}
	return nil  // for 1.0
}

func (p Aio) Extract(retchan chan record.Record) {
	filename := "/proc/sys/fs/aio-nr"
	s, err := ReadLines(filename)
	if err != io.EOF {
		close(retchan)
		return
	}
	r := record.Record{p.Name, time.Now(), map[string]string{
		p.Name: s[0],
	}}
	retchan <- r
	return
}
