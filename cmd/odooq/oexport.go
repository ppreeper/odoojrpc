package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/odoojrpc"
	"github.com/urfave/cli/v2"
)

var Export = cli.Command{
	Name:        "export",
	UsageText:   "export - export model records",
	Description: "exports records from odoo via rpc call",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "model", Aliases: []string{"m"}, Destination: &Model, Usage: "`MODEL` to export", Required: true},
		&cli.StringFlag{Name: "field", Aliases: []string{"f"}, Destination: &Field, Usage: "fields to export `FIELD`"},
		&cli.StringFlag{Name: "domain", Aliases: []string{"d"}, Destination: &Domain, Usage: "filter `DOMAIN`"},
		&cli.IntFlag{Name: "offset", Aliases: []string{"o"}, Destination: &RecordOffset, Value: 0, Usage: "offset `OFFSET` records from beginning"},
		&cli.IntFlag{Name: "limit", Aliases: []string{"l"}, Destination: &RecordLimit, Value: 0, Usage: "limit `LIMIT` records returned"},
		&cli.StringFlag{Name: "file", Destination: &DataFile, Usage: "export file `FILE`, if blank uses `MODEL` as basename"},
		&cli.StringFlag{Name: "separator", Aliases: []string{"s"}, Value: ";", Destination: &FileSep, Usage: "CSV separator"},
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
		// get field names
		fields := mapkeys(recs)

		// get filename
		if DataFile == "" {
			DataFile = c.String("model") + ".csv"
		}

		// csv writer setup
		fexport, err := os.Create(DataFile)
		if err != nil {
			fmt.Println(err)
		}
		defer fexport.Close()

		writer := csv.NewWriter(fexport)
		if FileSep != "" {
			writer.Comma = []rune(FileSep)[0]
		}

		// header
		err = writer.Write(fields)
		if err != nil {
			fmt.Println(err)
		}
		writer.Flush()

		// records
		for _, rec := range recs {
			v := maptoslice(rec, fields)
			err = writer.Write(v)
			if err != nil {
				fmt.Println(err)
			}
			writer.Flush()
		}

		return nil
	},
}
