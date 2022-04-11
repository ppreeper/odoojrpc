package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ppreeper/odoojrpc"
)

var O *OdooConn

// type alias to reduce typing
type oarg = odoojrpc.FilterArg

// OdooConn structure to provide basic connection
type OdooConn struct {
	Hostname string `default:"localhost"`
	Port     int    `default:"8069"`
	Database string `default:"odoo"`
	Username string `default:"odoo"`
	Password string `default:"odoo"`
	Schema   string `default:"http"`
	Protocol string `default:"xmlrpcs"`
	Workers  int    `default:"1"`
	*odoojrpc.Odoo
}

//NewOdooConn initializer
func NewOdooConn(oc OdooConn) *OdooConn {
	oc.Odoo = &odoojrpc.Odoo{
		Hostname: oc.Hostname,
		Port:     oc.Port,
		Username: oc.Username,
		Password: oc.Password,
		Schema:   oc.Schema,
		Database: oc.Database,
	}
	return &oc
}

func ModelName(mdl string) string {
	return strings.Replace(mdl, "_", ".", -1)
}

func login() {
	O = NewOdooConn(OdooConn{Hostname: hostname, Database: database, Username: username, Password: password, Protocol: protocol, Schema: schema, Port: port, Workers: workers})
	err := O.Login()
	if err != nil {
		fmt.Printf("login error: %v\n", err)
		os.Exit(1)
	}
}
