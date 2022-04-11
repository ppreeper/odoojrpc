package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Export = cli.Command{
	Name: "export",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "model", Usage: "`MODEL` to export", Required: true},
		&cli.StringFlag{Name: "field", Usage: "fields to export `FIELD`", Required: true},
		&cli.StringFlag{Name: "domain", Usage: "filter `DOMAIN`"},
		&cli.StringFlag{Name: "file", Usage: "output file `FILE`", Required: true},
		&cli.IntFlag{Name: "size", Value: 10, Usage: "Number of lines to import per connection `SIZE`"},
		&cli.StringFlag{Name: "separator", Aliases: []string{"s"}, Value: ";", Usage: "CSV separator"},
	},
	Action: func(c *cli.Context) error {
		fmt.Println("oexport")

		fmt.Println("args.len", c.Args().Len())

		fmt.Println("model:", c.String("model"))
		fmt.Println("model:", c.String("field"))
		fmt.Println("model:", c.String("domain"))
		fmt.Println("model:", c.String("file"))
		fmt.Println("model:", c.String("size"))
		fmt.Println("model:", c.String("separator"))

		return nil
	},
}
