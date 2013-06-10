package modules

import (
	"bitbucket.org/r_rudi/gostat/record"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Plugin interface {
	Check(map[string]map[string]string) error                 // check that resource is available or not
	Extract(chan record.Record, map[string]map[string]string) // get information from that resource

}

// Read contents from file and split by new line.
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

type ModuleError struct {
	ModuleName string
	Where      string
	Reason     string
}

func NewModuleError(modulename string, where string, reason string) error {
	return ModuleError{
		modulename,
		where,
		reason,
	}
}
func (e ModuleError) Error() string {
	return fmt.Sprintf("%s: %s, %s", e.ModuleName, e.Where, e.Reason)
}
