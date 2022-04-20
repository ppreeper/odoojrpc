package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/ppreeper/odoojrpc"
	"github.com/urfave/cli/v2"
)

var Import = cli.Command{
	Name:        "import",
	UsageText:   "import - import records to a model",
	Description: "imports records to odoo via rpc call",
	Flags: []cli.Flag{

		&cli.StringFlag{Name: "model", Aliases: []string{"m"}, Destination: &Model, Usage: "`MODEL` to import", Required: true},
		&cli.StringFlag{Name: "field", Aliases: []string{"f"}, Destination: &Field, Usage: "fields to import `FIELD`", Required: true},
		&cli.IntFlag{Name: "offset", Aliases: []string{"o"}, Destination: &RecordOffset, Value: 0, Usage: "offset `OFFSET` records from beginning"},
		&cli.IntFlag{Name: "limit", Aliases: []string{"l"}, Destination: &RecordLimit, Value: 0, Usage: "limit `LIMIT` records returned"},
		&cli.StringFlag{Name: "file", Destination: &DataFile, Usage: "import file `FILE`, if blank uses MODEL as basename"},
		// &cli.IntFlag{Name: "size", Value: 10, Destination: &BatchSize, Usage: "Number of lines to import per connection `SIZE`"},
		&cli.StringFlag{Name: "separator", Aliases: []string{"s"}, Value: ";", Destination: &FileSep, Usage: "CSV separator"},

		// &cli.StringFlag{Name: "groupby", Aliases: []string{"g"}, Usage: "Group data per batch with the same value for the given column in order to avoid concurrent update error"},
		// &cli.StringFlag{Name: "ignore", Aliases: []string{"i"}, Destination: &FieldIgnore, Usage: "list of columns to `ignore` separated by comma. Those column will be remove from the import request"},
		// &cli.BoolFlag{Name: "check", Aliases: []string{"c"}, Value: true, Usage: "Check if record are imported after each batch"},
		// &cli.StringFlag{Name: "context", Usage: "`context` that will be passed to the load function, need to be a valid python dict"},
		// &cli.BoolFlag{Name: "o2m", Value: false, Usage: "When you want to import o2m field, don't cut the batch until we find a new id"},
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

		// get filename
		ErrorFile := c.String("file") + "_err.csv"
		if DataFile == "" {
			DataFile = c.String("model") + ".csv"
			ErrorFile = c.String("model") + "_err.csv"
		}

		fimport, err := os.Open(DataFile)
		if err != nil {
			fmt.Println("file open error", err)
			os.Exit(1)
		}
		defer fimport.Close()

		ferror, err := os.Create(ErrorFile)
		if err != nil {
			fmt.Println("file open error", err)
			os.Exit(1)
		}
		defer ferror.Close()

		// check record count
		linecount, err := lineCount(DataFile)
		if err != nil {
			fmt.Println("file open error", err)
			os.Exit(1)
		}
		if linecount <= 1 {
			fmt.Println("file too short")
			return nil
		}
		if RecordOffset > linecount+1 {
			fmt.Println("offset too large")
			return nil
		}
		// rewind file pointer to beginning of file
		fimport.Seek(0, io.SeekStart)

		// csv reader setup
		reader := csv.NewReader(fimport)
		if FileSep != "" {
			reader.Comma = []rune(FileSep)[0]
		}

		// get field names
		fields, err := reader.Read()
		if err != nil {
			fmt.Println("reader error:", err)
		}

		// csv error writer
		writer := csv.NewWriter(ferror)
		if FileSep != "" {
			writer.Comma = []rune(FileSep)[0]
		}

		err = writer.Write(fields)
		if err != nil {
			fmt.Println(err)
		}
		writer.Flush()
		offset := c.Int("offset")
		lcount := 1
		for {
			rec, err := reader.Read()
			if err == io.EOF {
				break
			}
			if offset > 0 {
				offset--
				continue
			}
			ur, err := slicetomap(fields, rec)
			if err != nil {
				fmt.Println(err)
			}
			id, err := strconv.Atoi(ur["id"].(string))
			if err != nil {
				fmt.Println(err)
			}
			umdl := ModelName(Model)
			if id == -1 {
				_, err := o.Create(ModelName(Model), ur)
				if err != nil {
					fmt.Println(err)
					writer.Write(rec)
					writer.Flush()
				}
			} else {
				_, err := o.Update(umdl, id, ur)
				if err != nil {
					fmt.Println(err)
					writer.Write(rec)
					writer.Flush()
				}
			}
			if lcount >= c.Int("limit") {
				break
			}
			lcount++
		}
		return nil
	},
}
