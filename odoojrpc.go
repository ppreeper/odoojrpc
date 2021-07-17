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

type FilterArg []interface{}

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
func (o *Odoo) Call(service string, method string, args ...interface{}) (res interface{}, err error) {
	params := map[string]interface{}{
		"service": service,
		"method":  method,
		"args":    args,
	}
	res, err = o.JSONRPC(params)
	return
}

// JSONRPC json request
func (o *Odoo) JSONRPC(params map[string]interface{}) (out interface{}, err error) {
	rand.Seed(time.Now().UnixNano())

	message := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "call",
		"id":      rand.Intn(100000000),
		"params":  params,
	}
	// fmt.Printf("\n\n message: %v", message)
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("json marshall error: %w", err)
	}
	resp, err := http.Post(o.URL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		err = fmt.Errorf("http post error: %w", err)
	}

	var result map[string]interface{}
	if resp != nil {
		json.NewDecoder(resp.Body).Decode(&result)
	}

	// fmt.Printf("\n\n resErr: %v", result["error"])
	if resErr, ok := result["error"]; ok {
		return resErr, err
	}
	// fmt.Printf("\n\n resVal: %v", result["result"])
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
	case interface{}:
		fmt.Println(v)
	}
	return
}

// Create record
func (o *Odoo) Create(model string, record map[string]interface{}) (out int, err error) {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "create", record)
	if err != nil {
		return -1, err
	}
	// fmt.Printf("\n\n Create: %v", v)
	switch v := v.(type) {
	case float64:
		out = int(v)
	case interface{}:
		// code := int(v.(map[string]interface{})["code"].(float64))
		svrMessage := v.(map[string]interface{})["message"].(string)
		data := v.(map[string]interface{})["data"].(map[string]interface{})
		name := data["name"].(string)
		message := data["message"].(string)
		err = fmt.Errorf("create record error model: %s message: %s %s %s record: %v", model, svrMessage, name, message, record)
		out = -1
	default:
		out = -1
	}
	return
}

// SearchRead records
func (o *Odoo) SearchRead(model string, filter FilterArg, offset int, limit int, fields []string) []map[string]interface{} {
	var dd []map[string]interface{}
	vv, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "search_read", filter, fields, offset, limit)
	if err != nil {
		log.Println(err)
	}
	switch vv := vv.(type) {
	case []interface{}:
		for _, v := range vv {
			dd = append(dd, v.(map[string]interface{}))
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
	case []interface{}:
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
	case []interface{}:
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
func (o *Odoo) Read(model string, ids []int, fields []string) []map[string]interface{} {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "read", ids, fields)
	if err != nil {
		log.Println(err)
	}
	var dd []map[string]interface{}
	switch v := v.(type) {
	case []interface{}:
		for _, v := range v {
			dd = append(dd, v.(map[string]interface{}))
		}
	}
	return dd
}

// Update record
func (o *Odoo) Update(model string, recordID int, record map[string]interface{}) (out bool, err error) {
	v, err := o.Call("object", "execute", o.Database, o.uid, o.Password, model, "write", recordID, record)
	if err != nil {
		return false, err
	}
	switch v := v.(type) {
	case bool:
		out = v
	case interface{}:
		// code := int(v.(map[string]interface{})["code"].(float64))
		svrMessage := v.(map[string]interface{})["message"].(string)
		data := v.(map[string]interface{})["data"].(map[string]interface{})
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
		filter = []interface{}{FilterArg{"id", "!=", "-1"}}
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

// Common Odoo Queries

// CompanyID record
func (o *Odoo) CompanyID(companyName string) int {
	return o.GetID("res.company", []interface{}{FilterArg{"name", "=", companyName}})
}

// PartnerID record
func (o *Odoo) PartnerID(partnerName string) int {
	return o.GetID("res.partner", []interface{}{FilterArg{"name", "=", partnerName}})
}

// CountryID record
func (o *Odoo) CountryID(countryName string) int {
	return o.GetID("res.country", []interface{}{FilterArg{"name", "=", countryName}})
}

// StateID record
func (o *Odoo) StateID(countryID int, stateName string) int {
	return o.GetID("res.country.state", []interface{}{FilterArg{"name", "=", stateName}, FilterArg{"country_id", "=", countryID}})
}

// FiscalPosition record
func (o *Odoo) FiscalPosition(countryID int, fiscalName string) int {
	return o.GetID("account.fiscal.position", []interface{}{FilterArg{"country_id", "=", countryID}, FilterArg{"name", "=", fiscalName}})
}
