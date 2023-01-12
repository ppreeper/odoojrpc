package odoojrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	uid      int
}

type FilterArg []any

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
	return
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
	return
}

// Call sends a request
func (o *Odoo) Call(service string, method string, args ...any) (res any, err error) {
	params := map[string]any{
		"service": service,
		"method":  method,
		"args":    args,
	}
	res, err = o.JSONRPC(params)
	return
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
		err = fmt.Errorf("http post error: %w", err)
	}

	var result map[string]any
	if resp != nil {
		json.NewDecoder(resp.Body).Decode(&result)
	}

	if _, ok := result["error"]; ok {
		resError := ""
		if errorMessage, ok := result["error"].(map[string]any)["message"].(string); ok {
			resError += errorMessage
		}
		if dataMessage, ok := result["error"].(map[string]any)["data"].(map[string]any)["message"].(string); ok {
			resError += ": " + dataMessage
		}
		return resError, err
	}
	resVal := result["result"]
	return resVal, err
}

// Login connects to server
func (o *Odoo) Login() (err error) {
	if o.URL == "" {
		err = o.Init()
		if err != nil {
			return
		}
	}
	v, err := o.Call("common", "login", o.Database, o.Username, o.Password)
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}
	switch v := v.(type) {
	case float64:
		o.uid = int(v)
	case string:
		return fmt.Errorf(v)
	}
	return
}

// Create record
func (o *Odoo) Create(model string, record map[string]any) (out int, err error) {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "create", record)
	if err != nil {
		return -1, err
	}
	// fmt.Printf("\n\n Create: %v", v)
	switch v := v.(type) {
	case float64:
		out = int(v)
	case any:
		// code := int(v.(map[string]any)["code"].(float64))
		svrMessage := v.(map[string]any)["message"].(string)
		data := v.(map[string]any)["data"].(map[string]any)
		name := data["name"].(string)
		message := data["message"].(string)
		err = fmt.Errorf("create record error model: %s message: %s %s %s record: %v", model, svrMessage, name, message, record)
		out = -1
	default:
		out = -1
	}
	return
}

// Load record
func (o *Odoo) Load(model string, header []string, records []any) (out int, err error) {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "load", header, records)
	if err != nil {
		return -1, err
	}
	switch v := v.(type) {
	case float64:
		out = int(v)
	case map[string]any:
		if v["message"] != nil {
			err = fmt.Errorf("create record error model: %s ids %v message %v", model, v["ids"], v["message"])
		}
		out = 0
	default:
		out = -1
	}
	return
}

// SearchRead records
func (o *Odoo) SearchRead(model string, filter FilterArg, offset int, limit int, fields []string) []map[string]any {
	var dd []map[string]any
	vv, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "search_read", filter, fields, offset, limit)
	if err != nil {
		log.Println(err)
	}
	switch vv := vv.(type) {
	case []any:
		for _, v := range vv {
			dd = append(dd, v.(map[string]any))
		}
	}
	return dd
}

// Search record
func (o *Odoo) Search(model string, filter FilterArg) (oo []int) {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "search", filter)
	if err != nil {
		log.Println(err)
	}
	switch v := v.(type) {
	case []any:
		for _, v := range v {
			oo = append(oo, int(v.(float64)))
		}
	}
	return
}

// GetID record
func (o *Odoo) GetID(model string, filter FilterArg) (out int) {
	out = -1
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "search", filter)
	if err != nil {
		log.Println(err)
	}
	switch v := v.(type) {
	case []any:
		dd := []int{}
		for _, v := range v {
			dd = append(dd, int(v.(float64)))
		}
		if len(dd) > 0 {
			out = dd[0]
		}
	}
	return
}

// Read record
func (o *Odoo) Read(model string, ids []int, fields []string) []map[string]any {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "read", ids, fields)
	if err != nil {
		log.Println(err)
	}
	var dd []map[string]any
	switch v := v.(type) {
	case []any:
		for _, v := range v {
			dd = append(dd, v.(map[string]any))
		}
	}
	return dd
}

// Update record
func (o *Odoo) Update(model string, recordID int, record map[string]any) (out bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "write", recordID, record)
	if err != nil {
		return false, err
	}
	switch v := v.(type) {
	case bool:
		out = v
	case any:
		svrMessage := v.(map[string]any)["message"].(string)
		data := v.(map[string]any)["data"].(map[string]any)
		name := data["name"].(string)
		message := data["message"].(string)
		err = fmt.Errorf("update record error model: %s message: %s %s %s record: %v", model, svrMessage, name, message, record)
		out = false
	default:
		out = false
	}
	return
}

// Unlink record
func (o *Odoo) Unlink(model string, recordIDs []int) (out bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "unlink", recordIDs)
	if err != nil {
		log.Println(err)
	}
	switch v := v.(type) {
	case bool:
		out = v
	default:
		out = false
	}
	return
}

// Count record
func (o *Odoo) Count(model string, filter FilterArg) (out int) {
	if len(filter) == 0 {
		filter = []any{FilterArg{"id", "!=", "-1"}}
	}
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "search_count", filter)
	if err != nil {
		log.Println(err)
	}
	switch v := v.(type) {
	case float64:
		out = int(v)
	default:
		out = -1
	}
	return
}
