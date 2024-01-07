// odoojrpc - go library to access Odoo server via Json RPC
// Copyright (C) 2021  Peter Preeper

// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.

// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.

// You should have received a copy of the GNU Lesser General Public
// License along with this library; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301
// USA
package odoojrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

// Odoo connection
// Return a new instance of the :class 'Odoo' class.
type Odoo struct {
	Hostname string `default:"localhost"`
	Port     int    `default:"8069"`
	Database string `default:"odoo"`
	Username string `default:"odoo"`
	Password string `default:"odoo"`
	Schema   string `default:"http"`
	URL      string
	UID      int
}

func (o Odoo) WithHostname(hostname string) Odoo {
	o.Hostname = hostname
	return o
}

func (o Odoo) WithPort(port int) Odoo {
	o.Port = port
	return o
}

func (o Odoo) WithDatabase(database string) Odoo {
	o.Database = database
	return o
}

func (o Odoo) WithUsername(username string) Odoo {
	o.Username = username
	return o
}

func (o Odoo) WithPassword(password string) Odoo {
	o.Password = password
	return o
}

func (o Odoo) WithSchema(schema string) Odoo {
	o.Schema = schema
	return o
}

func NewOdoo() *Odoo {
	return &Odoo{}
}

func NewOdooWithConfig(hostname string, port int, database string, username string, password string, schema string) *Odoo {
	return &Odoo{
		Hostname: hostname,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
		Schema:   schema,
	}
}

var (
	// ErrLogin error on login failure
	ErrLogin   = errors.New("login failed")
	ErrSchema  = errors.New("invalid schema: http or https")
	ErrPort    = errors.New("invalid port: 1-65535")
	ErrHostLen = errors.New("invalid hostname length: 1-2048")
)

func (o *Odoo) Init() (err error) {
	if err = o.genURL(); err != nil {
		return fmt.Errorf("init error: %w", err)
	}
	return nil
}

// genURL returns url string
func (o *Odoo) genURL() (err error) {
	if o.Schema != "http" && o.Schema != "https" {
		err = ErrSchema
		return
	}
	if o.Port == 0 || o.Port > 65535 {
		err = ErrPort
		return
	}
	if len(o.Hostname) == 0 || len(o.Hostname) > 2048 {
		err = ErrHostLen
		return
	}

	o.URL = fmt.Sprintf("%s://%s:%d/jsonrpc", o.Schema, o.Hostname, o.Port)
	return nil
}

// Call sends a request
func (o *Odoo) Call(service string, method string, args ...any) (res any, err error) {
	params := map[string]any{
		"service": service,
		"method":  method,
		"args":    args,
	}
	res, err = o.JSONRPC(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// JSONRPC json request
func (o *Odoo) JSONRPC(params map[string]any) (out any, err error) {
	message := map[string]any{
		"jsonrpc": "2.0",
		"method":  "call",
		"id":      rand.Intn(100000000),
		"params":  params,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("json marshall error: %w", err)
	}

	resp, err := http.Post(o.URL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, fmt.Errorf("http post error: %w", err)
	}

	var result map[string]any
	if resp != nil {
		json.NewDecoder(resp.Body).Decode(&result)
	} else {
		return nil, fmt.Errorf("no response returned")
	}

	if _, ok := result["error"]; ok {
		resultError := ""
		if errorMessage, ok := result["error"].(map[string]any)["message"].(string); ok {
			resultError += errorMessage
		}
		if dataMessage, ok := result["error"].(map[string]any)["data"].(map[string]any)["message"].(string); ok {
			resultError += ": " + dataMessage
		}
		return nil, fmt.Errorf(resultError)
	}

	out = result["result"]
	return out, nil
}

// Login connects to server
func (o *Odoo) Login() (err error) {
	if o.URL == "" {
		err = o.Init()
		if err != nil {
			return err
		}
	}
	v, err := o.Call("common", "login", o.Database, o.Username, o.Password)
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}
	switch v := v.(type) {
	case float64:
		o.UID = int(v)
	}
	return nil
}

// Create record
func (o *Odoo) Create(model string, record map[string]any) (row int, res bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "create", record)
	if err != nil {
		return -1, false, err
	}
	switch v := v.(type) {
	case float64:
		row = int(v)
	default:
		row = -1
	}
	return row, res, nil
}

// Load record
func (o *Odoo) Load(model string, header []string, records []any) (row int, res bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "load", header, records)
	if err != nil {
		return -1, false, err
	}
	switch v := v.(type) {
	case float64:
		row = int(v)
	default:
		row = -1
	}
	return row, res, nil
}

// SearchRead records
func (o *Odoo) SearchRead(model string, filter []any, offset int, limit int, fields []string) (oo []map[string]any, err error) {
	vv, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "search_read", filter, fields, offset, limit)
	if err != nil {
		return oo, err
	}
	switch vv := vv.(type) {
	case []any:
		for _, v := range vv {
			oo = append(oo, v.(map[string]any))
		}
	}
	return oo, nil
}

// Search record
func (o *Odoo) Search(model string, filter []any) (rows []int, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "search", filter)
	if err != nil {
		return rows, err
	}
	switch v := v.(type) {
	case []any:
		for _, v := range v {
			rows = append(rows, int(v.(float64)))
		}
	}
	return rows, nil
}

// GetID record
func (o *Odoo) GetID(model string, filter []any) (out int, res bool, err error) {
	out = -1
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "search", filter)
	if err != nil {
		return out, false, err
	}
	switch v := v.(type) {
	case []any:
		rr := []int{}
		for _, v := range v {
			rr = append(rr, int(v.(float64)))
		}
		if len(rr) > 0 {
			out = rr[0]
		}
	}
	return out, res, nil
}

// Read record
func (o *Odoo) Read(model string, ids []int, fields []string) (oo []map[string]any, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "read", ids, fields)
	if err != nil {
		return oo, err
	}
	switch v := v.(type) {
	case []any:
		for _, v := range v {
			oo = append(oo, v.(map[string]any))
		}
	}
	return oo, nil
}

// Update record
func (o *Odoo) Update(model string, recordID int, record map[string]any) (row int, res bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "write", recordID, record)
	if err != nil {
		return recordID, false, err
	}
	switch v := v.(type) {
	case bool:
		res = v
	default:
		res = false
	}
	return recordID, res, nil
}

// Unlink record
func (o *Odoo) Unlink(model string, recordIDs []int) (res bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "unlink", recordIDs)
	if err != nil {
		return res, err
	}
	switch v := v.(type) {
	case bool:
		res = v
	default:
		res = false
	}
	return res, nil
}

// Count record
func (o *Odoo) Count(model string, filter []any) (count int, err error) {
	if len(filter) == 0 {
		filter = []any{[]any{"id", "!=", "-1"}}
	}
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "search_count", filter)
	if err != nil {
		return count, err
	}
	switch v := v.(type) {
	case float64:
		count = int(v)
	default:
		count = -1
	}
	return count, nil
}
