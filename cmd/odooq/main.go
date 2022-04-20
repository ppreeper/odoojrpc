package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/ppreeper/odoojrpc"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// type alias to reduce typing
type oarg = odoojrpc.FilterArg

// connection variables
var (
	cfgFile  string
	hostname string
	database string
	username string
	password string
	protocol string
	schema   string
	port     int
	workers  int
)

// common variables
var (
	Model        string
	Field        string
	FieldIgnore  string
	Fields       []string
	Domain       string
	RecordOffset = 0
	RecordLimit  = 0
	BatchSize    = 1
	DataFile     string
	FileSep      string
)

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

func ModelName(mdl string) string {
	return strings.Replace(mdl, "_", ".", -1)
}

func lineCount(filename string) (int, error) {
	lc := int(0)
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		lc++
	}
	return lc, s.Err()
}

func maptoslice(record map[string]interface{}, fields []string) (values []string) {
	for _, key := range fields {
		values = append(values, fmt.Sprint(record[key]))
	}
	return
}

func mapkeys(records []map[string]interface{}) (keys []string) {
	for k := range records[0] {
		keys = append(keys, k)
	}
	index := -1
	for k, v := range keys {
		if v == "id" {
			index = k
			break
		}
	}
	copy(keys[index:], keys[index+1:])
	keys[len(keys)-1] = ""
	keys = keys[:len(keys)-1]
	keys = append([]string{"id"}, keys...)
	return
}

func slicetomap(keys []string, values []string) (map[string]interface{}, error) {
	var v = make(map[string]interface{})
	if len(keys) != len(values) {
		return v, fmt.Errorf("incorrect key value count")
	}
	for i := 0; i < len(keys); i++ {
		v[keys[i]] = values[i]
	}
	return v, nil
}
