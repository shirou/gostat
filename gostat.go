package main

import (
	"bitbucket.org/r_rudi/gostat/modules"
	"bitbucket.org/r_rudi/gostat/outputs"
	"bitbucket.org/r_rudi/gostat/record"
	"flag"
	"github.com/msbranco/goconfig"
	"reflect"
	"runtime"
	"time"
	"sync"
)

func Call(m map[string]func() (modules.Plugin, error), name string) (result interface{}, err error) {
	f := reflect.ValueOf(m[name])
	return f.Call(nil), nil
}

func funcs() map[string]func() modules.Plugin {
	funcs := map[string]func() modules.Plugin{
		"aio":  modules.NewAio,
		"cpu":  modules.NewCpu,
		"load": modules.NewLoad,
		"mem":  modules.NewMem,
	}
	return funcs
}

func get(plugin_list []modules.Plugin, out outputs.Output, conf map[string]map[string]string) {
	ch := make(chan record.Record)
	ret_list := make([]record.Record, 0)
	var wg sync.WaitGroup

	wg.Add(len(plugin_list))
	for _, o := range plugin_list {
		go func(){
			defer wg.Done()
			o.Extract(ch, conf)
		}()
	}

	go func(){
		defer wg.Done()
		for msg := range ch {
			ret_list = append(ret_list, msg)
		}
	}()

	wg.Wait()

	out.Emit(ret_list, conf)
}

func main() {
	c := flag.String("c", "~/.gostat.conf", "Config file")
	o := flag.String("o", "ltsv", "Output format")
	i := flag.Int("i", 0, "interval time(seconds)")
	flag.Parse()

	// load config file
	conf := make(map[string]map[string]string)
	conf["root"] = make(map[string]string)
	conf["root"]["os"] = runtime.GOOS
	conf["root"]["configfile"] = *c
	if *c != "" {
		config, err := goconfig.ReadConfigFile(*c)
		if err == nil {
			for _, section := range config.GetSections() {
				options, _ := config.GetOptions(section)
				for _, option := range options {
					if _, ok := conf[section]; ok == false {
						conf[section] = make(map[string]string)
					}
					conf[section][option], _ = config.GetRawString(section, option)
				}
			}
		}
	}

	plugin_list := make([]modules.Plugin, 0)

	f := funcs()
	keys := make([]string, 0)
	for k := range f {
		keys = append(keys, k)
	}

	for _, plugin := range keys {
		ret := reflect.ValueOf(f[plugin]).Call(nil)
		module := ret[0].Interface().(modules.Plugin)
		if module.Check(conf) == nil {
			plugin_list = append(plugin_list, module)
		}
	}

	// Setting Output format
	var out outputs.Output
	switch *o {
	case "ltsv":
		out = outputs.Output(outputs.LTSV{})
	case "whitespace":
		out = outputs.Output(outputs.WhiteSpace{})
	case "csv":
		out = outputs.Output(outputs.CSV{})
	case "mqtt":
		out = outputs.Output(outputs.MQTT{})
	case "http":
		out = outputs.Output(outputs.Http{})
	default:
		panic("Unknown output format")
	}

	// Extract
	if *i == 0 {
		get(plugin_list, out, conf)
	} else {
		for {
			get(plugin_list, out, conf)
			time.Sleep(time.Second * time.Duration(*i))
		}
	}
}
