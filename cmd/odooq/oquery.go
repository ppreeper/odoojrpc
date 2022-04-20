package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/odoojrpc"
	"github.com/urfave/cli/v2"
)

var Query = cli.Command{
	Name:        "query",
	UsageText:   "query - query model for records",
	Description: "queries records from odoo via rpc call",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "model", Aliases: []string{"m"}, Destination: &Model, Usage: "`model` to query", Required: true},
		&cli.StringFlag{Name: "field", Aliases: []string{"f"}, Destination: &Field, Usage: "fields to export `FIELD`"},
		&cli.StringFlag{Name: "domain", Aliases: []string{"d"}, Destination: &Domain, Usage: "filter `DOMAIN`"},
		&cli.IntFlag{Name: "offset", Aliases: []string{"o"}, Destination: &RecordOffset, Value: 0, Usage: "offset `OFFSET` records from beginning"},
		&cli.IntFlag{Name: "limit", Aliases: []string{"l"}, Destination: &RecordLimit, Value: 0, Usage: "limit `LIMIT` records returned"},
	},
	Action: func(c *cli.Context) error {
		var o = &odoojrpc.Odoo{
			Hostname: hostname,
			Port:     port,
			Username: username,
			Password: password,
			Schema:   schema,
			Database: database,
		}

		err := o.Login()
		if err != nil {
			fmt.Println("login error", err)
			os.Exit(1)
		}

		umdl := ModelName(Model)
		if Field != "" {
			Fields = strings.Split(Field, ",")
		}
		filter, err := odoojrpc.SearchDomain(Domain)
		if err != nil {
			fmt.Println("invalid domain:", err)
			os.Exit(1)
		}

		recs := o.SearchRead(umdl, filter, RecordOffset, RecordLimit, Fields)
		if len(recs) <= 0 {
			fmt.Println("no records found")
			return nil
		}
		j, err := json.MarshalIndent(recs, "", "  ")
		if err != nil {
			fmt.Println("error processing records", err)
			os.Exit(1)
		}
		fmt.Println(string(j))

		return nil
	},
}
