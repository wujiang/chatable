package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/glog"
	"github.com/wujiang/chatable/api"
	"github.com/wujiang/chatable/datastore"
	"github.com/wujiang/chatable/rds"
)

var (
	baseURL    *url.URL
	cfg        = &configuration{}
	configFile = flag.String("c", "development.json", "Configuration file")
)

type subCMD struct {
	name string
	desc string
	exec func(args []string)
}

var subcmds = []subCMD{
	{name: "serve", desc: "serve the server", exec: serveCMD},
}

func main() {
	// command usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
chatable

Usage:
        chatable [options] command [arg...]

Commands
`)
		for _, cmd := range subcmds {
			fmt.Fprintf(os.Stderr, "\t%s - %s\n", cmd.name, cmd.desc)
		}
		fmt.Fprintf(os.Stderr, `
Use "chatable command -h" for command help.

Options:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}

	// parse the configuration file
	err := cfg.parse(*configFile)
	if err != nil {
		glog.Fatal(err)
	}

	// initialize database connection
	datastore.Init(cfg.Postgres)
	defer datastore.Exit()

	// initialize redis
	rds.Init(cfg.RedisHost)

	rdsPool := rds.NewRdsPool(nil)
	rdsPool.AddToQM(cfg.QueueManagerKey, cfg.MessageQueueKey)
	defer rdsPool.RemoveFromQM(cfg.QueueManagerKey, cfg.MessageQueueKey)

	go api.Hub.Run(cfg.SharedQueueKey)

	qm := api.QueueManager{}
	go qm.Dispatch(cfg.SharedQueueKey, cfg.QueueManagerKey)
	go qm.Pop(cfg.MessageQueueKey)

	subcmd := flag.Arg(0)
	for _, c := range subcmds {
		if c.name == subcmd {
			c.exec(flag.Args()[1:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command %q\n", subcmd)
	flag.Usage()
	os.Exit(1)
}

func serveCMD(args []string) {
	sv := flag.NewFlagSet("serve", flag.ExitOnError)
	sv.Usage = func() {
		fmt.Fprintf(os.Stderr, `
usage: chatable serve [options]

Start web server and api.

Options:
`)
		sv.PrintDefaults()
		os.Exit(1)
	}
	sv.Parse(args)
	if sv.NArg() != 0 {
		sv.Usage()
	}

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))

	fmt.Println("Serving on", cfg.Host)
	err := http.ListenAndServe(cfg.Host, m)
	if err != nil {
		glog.Fatal(err)
	}
}
