package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Import = cli.Command{
	Name: "import",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "model", Usage: "`MODEL` to export", Required: true},
		&cli.StringFlag{Name: "file", Usage: "output file `FILE`", Required: true},
		&cli.IntFlag{Name: "size", Value: 10, Usage: "Number of lines to import per connection `SIZE`"},
		&cli.IntFlag{Name: "skip", Usage: "Skip until line [`SKIP`]"},
		&cli.StringFlag{Name: "separator", Aliases: []string{"s"}, Value: ";", Usage: "CSV separator"},
		&cli.StringFlag{Name: "groupby", Usage: "Group data per batch with the same value for the given column in order to avoid concurrent update error"},
		&cli.StringFlag{Name: "ignore", Usage: "list of columns to `ignore` separated by comma. Those column will be remove from the import request"},
		&cli.BoolFlag{Name: "check", Value: true, Usage: "Check if record are imported after each batch"},
		&cli.StringFlag{Name: "context", Usage: "`context` that will be passed to the load function, need to be a valid python dict"},
		&cli.BoolFlag{Name: "o2m", Value: false, Usage: "When you want to import o2m field, don't cut the batch until we find a new id"},
	},
	Action: func(c *cli.Context) error {
		fmt.Println("oimport")

		fmt.Println("model:", c.String("model"))
		fmt.Println("file:", c.String("file"))
		fmt.Println("size:", c.String("size"))
		fmt.Println("skip:", c.String("skip"))
		fmt.Println("separator:", c.String("separator"))
		fmt.Println("groupby:", c.String("groupby"))
		fmt.Println("ignore:", c.String("ignore"))
		fmt.Println("check:", c.String("check"))
		fmt.Println("context:", c.String("context"))
		fmt.Println("o2m:", c.String("o2m"))

		return nil
	},
}
