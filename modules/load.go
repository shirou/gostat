package modules

import (
	"bitbucket.org/r_rudi/gostat/record"
	"io"
	"log"
	"os"
	"os/exec"
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

func NewLoad() Plugin {
	p := Load{
		"load",
		[]string{"1m", "5m", "15m"},
		[]string{"load1", "load5", "load15"},
		"f",
		0.5,
	}
	return Plugin(p)
}

func (p Load) Check(conf map[string]map[string]string) error {
	switch conf["root"]["os"] {
	case "linux":
		filename := "/proc/loadavg"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return err
		} else {
			return nil
		}
		return nil
	case "freebsd":
		return nil
	default:
		return NewModuleError(p.Name, "Check", "Unsupported platform")
	}
}

func (p Load) extractLinux() (map[string]string, error) {
	filename := "/proc/loadavg"
	s, err := ReadLines(filename)
	if err != io.EOF {
		return nil, err
	}

	values := strings.Fields(s[0])
	ret := map[string]string{}
	for i, t := range p.Vars {
		ret[t] = values[i]
	}
	return ret, nil
}

func (p Load) extractFreeBSD() (map[string]string, error) {
	out, err := exec.Command("/sbin/sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		log.Fatal(err)
	}
	v := strings.Replace(string(out), "{ ", "", 1)
	v = strings.Replace(string(v), " }", "", 1)
	values := strings.Fields(string(v))
	ret := map[string]string{}
	for i, t := range p.Vars {
		ret[t] = values[i]
	}
	return ret, err
}

func (p Load) Extract(retchan chan record.Record, conf map[string]map[string]string) {
	ret := map[string]string{}
	var err error
	switch conf["root"]["os"] {
	case "linux":
		ret, err = p.extractLinux()
	case "freebsd":
		ret, err = p.extractFreeBSD()
	default:
		return
	}
	if err != nil {
		return
	}

	r := record.Record{p.Name, time.Now(), ret}
	retchan <- r
	return
}
