package main

import (
	"./modules"
	"./outputs"
	"./record"
	"flag"
	"reflect"
	"time"
)

func Call(m map[string]func() (modules.Plugin, error), name string) (result interface{}, err error) {
	f := reflect.ValueOf(m[name])
	return f.Call(nil), nil
}

func funcs() map[string]func() (modules.Plugin, error) {
	funcs := map[string]func() (modules.Plugin, error){
		"aio":  modules.NewAio,
		"cpu":  modules.NewCpu,
		"load": modules.NewLoad,
		"mem":  modules.NewMem,
	}
	return funcs
}

func get(plugin_list []modules.Plugin, out outputs.Output, args []string) {
	ch := make(chan record.Record)
	count := 0
	for _, o := range plugin_list {
		count += 1
		go o.Extract(ch)
	}

	ret_list := make([]record.Record, 0)
	for i := 0; i < count; i++ {
		ret_list = append(ret_list, <-ch)
	}
	out.Emit(ret_list, args)

}

func main() {
	o := flag.String("o", "ltsv", "Output format")
	s := flag.Int("s", 0, "sleep seconds")
	flag.Parse()

	plugin_list := make([]modules.Plugin, 0)

	f := funcs()
	keys := make([]string, 0)
	for k := range f {
		keys = append(keys, k)
	}

	for _, plugin := range keys {
		ret := reflect.ValueOf(f[plugin]).Call(nil)
		if ret[1].IsNil() {
			plugin_list = append(plugin_list, ret[0].Interface().(modules.Plugin))
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
	if *s == 0 {
		get(plugin_list, out, flag.Args())
	} else {
		for {
			get(plugin_list, out, flag.Args())
			time.Sleep(time.Second * time.Duration(*s))
		}
	}

}
