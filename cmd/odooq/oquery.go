package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/odoojrpc"
	"github.com/urfave/cli/v2"
)

var (
	QueryModel  string
	QueryField  string
	QueryFields []string
	QueryDomain string
	QueryOffset = 0
	QueryLimit  = 0
)

var Query = cli.Command{
	Name:        "query",
	UsageText:   "query - query model for data",
	Description: "queries odoo via rpc call",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "model", Destination: &QueryModel, Usage: "`model` to query", Required: true},
		&cli.StringFlag{Name: "field", Destination: &QueryField, Usage: "fields to export `FIELD`"},
		&cli.StringFlag{Name: "domain", Destination: &QueryDomain, Usage: "filter `DOMAIN`"},
		&cli.IntFlag{Name: "offset", Destination: &QueryOffset, Value: 0, Usage: "offset `OFFSET` records from beginning"},
		&cli.IntFlag{Name: "limit", Destination: &QueryLimit, Value: 0, Usage: "limit `LIMIT` records returned"},
	},
	Action: func(c *cli.Context) error {
		umdl := ModelName(QueryModel)
		if QueryField != "" {
			QueryFields = strings.Split(QueryField, ",")
		}
		args, err := odoojrpc.SearchDomain(QueryDomain)
		if err != nil {
			fmt.Println("invalid domain:", err)
			os.Exit(1)
		}

		login()

		recs := O.SearchRead(umdl, args, QueryOffset, QueryLimit, QueryFields)
		j, err := json.MarshalIndent(recs, "", "  ")
		if err != nil {
			fmt.Println("error processing records", err)
			os.Exit(1)
		}
		fmt.Println(string(j))

		return nil
	},
}
