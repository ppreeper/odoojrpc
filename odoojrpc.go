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
	"time"
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

type (
	FilterArg []any
	oarg      = FilterArg
)

var (
	// ErrLogin error on login failure
	ErrLogin   = errors.New("login failed")
	ErrSchema  = errors.New("invalid schema: http or https")
	ErrPort    = errors.New("invalid port: 1-65535")
	ErrHostLen = errors.New("invalid hostname length: 1-2048")
)

func (o *Odoo) Init() (err error) {
	err = o.genURL()
	if err != nil {
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
	rand.Seed(time.Now().UnixNano())

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
func (o *Odoo) Create(model string, record map[string]any) (out int, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "create", record)
	if err != nil {
		return -1, err
	}
	switch v := v.(type) {
	case float64:
		out = int(v)
	default:
		out = -1
	}
	return out, nil
}

// Load record
func (o *Odoo) Load(model string, header []string, records []any) (out int, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "load", header, records)
	if err != nil {
		return -1, err
	}
	switch v := v.(type) {
	case float64:
		out = int(v)
	default:
		out = -1
	}
	return out, nil
}

// SearchRead records
func (o *Odoo) SearchRead(model string, filter FilterArg, offset int, limit int, fields []string) (oo []map[string]any, err error) {
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
func (o *Odoo) Search(model string, filter FilterArg) (oo []int, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "search", filter)
	if err != nil {
		return oo, err
	}
	switch v := v.(type) {
	case []any:
		for _, v := range v {
			oo = append(oo, int(v.(float64)))
		}
	}
	return oo, nil
}

// GetID record
func (o *Odoo) GetID(model string, filter FilterArg) (out int, err error) {
	out = -1
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "search", filter)
	if err != nil {
		return out, err
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
	return out, nil
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
func (o *Odoo) Update(model string, recordID int, record map[string]any) (out bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "write", recordID, record)
	if err != nil {
		return false, err
	}
	switch v := v.(type) {
	case bool:
		out = v
	default:
		out = false
	}
	return out, nil
}

// Unlink record
func (o *Odoo) Unlink(model string, recordIDs []int) (out bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "unlink", recordIDs)
	if err != nil {
		return out, err
	}
	switch v := v.(type) {
	case bool:
		out = v
	default:
		out = false
	}
	return out, nil
}

// Count record
func (o *Odoo) Count(model string, filter FilterArg) (out int, err error) {
	if len(filter) == 0 {
		filter = []any{FilterArg{"id", "!=", "-1"}}
	}
	v, err := o.Call("object", "execute", o.Database, o.UID, o.Password, model, "search_count", filter)
	if err != nil {
		return out, err
	}
	switch v := v.(type) {
	case float64:
		out = int(v)
	default:
		out = -1
	}
	return out, nil
}
