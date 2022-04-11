package main

import (
	"log"
	"os"
	"runtime"
	"sort"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var cfgFile string
var hostname string
var database string
var username string
var password string
var protocol string
var schema string
var port int
var workers int

func main() {
	var err error
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Destination: &cfgFile,
			Value:       "config.yml",
			FilePath:    ".",
			Usage:       "Load configuration from `FILE`",
			EnvVars:     []string{"ODOOQ_CONFIG"},
			HasBeenSet:  true,
		},
		altsrc.NewStringFlag(&cli.StringFlag{Name: "hostname", Destination: &hostname, Value: "localhost", Usage: "connect to `hostname`", EnvVars: []string{"ODOOQ_HOSTNAME"}}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "database", Destination: &database, Value: "odoo", Usage: "connect to `database`", EnvVars: []string{"ODOOQ_DATABASE"}}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "username", Destination: &username, Value: "admin", Usage: "login as `username`", EnvVars: []string{"ODOOQ_USERNAME"}}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "password", Destination: &password, Value: "admin", Usage: "using `password`", EnvVars: []string{"ODOOQ_PASSWORD"}}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "protocol", Destination: &protocol, Value: "jsonrpc", Usage: "connect using `protocol`", EnvVars: []string{"ODOOQ_PROTOCOL"}}),
		altsrc.NewStringFlag(&cli.StringFlag{Name: "schema", Destination: &schema, Value: "http", Usage: "connect using `protocol`", EnvVars: []string{"ODOOQ_SCHEMA"}}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "port", Value: 8069, Destination: &port, Usage: "connect to `port`", EnvVars: []string{"ODOOQ_PORT"}}),
		altsrc.NewIntFlag(&cli.IntFlag{Name: "workers", Value: runtime.NumCPU(), Destination: &workers, Usage: "number of simultaneous `workers`", EnvVars: []string{"ODOOQ_WORKERS"}}),
	}

	sort.Sort(cli.FlagsByName(flags))

	app := &cli.App{
		HelpName: "odooq",
		Flags:    flags,
		Before:   altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Commands: []*cli.Command{
			&Query,
			&Export,
			&Import,
		},
	}
	app.Setup()

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
