package odoojrpc

import (
	"reflect"
	"testing"
)

var searchDomainPatterns = []struct {
	domain string
	args   FilterArg
	err    error
}{
	{"", FilterArg{}, nil},
	{"('')", FilterArg{}, errSyntax},
	{"('','')", FilterArg{}, errSyntax},
	{"('a','=')", FilterArg{}, errSyntax},
	{"('name')", FilterArg{}, errSyntax},
	{"('name','=')", FilterArg{}, errSyntax},
	{"('name','=','My Name')", FilterArg{FilterArg{"name", "=", "My Name"}}, nil},
	{"('name','like','My Name')", FilterArg{FilterArg{"name", "like", "My Name"}}, nil},
	{"('name','=','My Name'),('name','=','My Name')", FilterArg{}, errSyntax},
	{"[('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}}, nil},
	{"[('name','=','My Name'),('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}}, nil},
	{"[('name','=','My Name'),'!',('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"!", FilterArg{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),'&',('name','=','My Name'),('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"&", FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"|", FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'!',('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}, FilterArg{"!", FilterArg{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'&',('name','=','My Name')]", FilterArg{}, errSyntax},
	{"[('name','=','My Name'),('name','=','My Name'),'&',('name','=','My Name'),('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}, FilterArg{"&", FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name')]", FilterArg{}, errSyntax},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}, FilterArg{"|", FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name'),'!',('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}, FilterArg{"|", FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}}, FilterArg{"!", FilterArg{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name'),'!',('name','=','My Name'),('name','=','My Name')]", FilterArg{FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}, FilterArg{"|", FilterArg{"name", "=", "My Name"}, FilterArg{"name", "=", "My Name"}}, FilterArg{"!", FilterArg{"name", "=", "My Name"}}, FilterArg{"name", "=", "My Name"}}, nil},
}

func TestSearchDomain(t *testing.T) {
	for i, pattern := range searchDomainPatterns {
		// fmt.Println("test: domain:", pattern.domain, "pattern.err:", pattern.err)
		args, err := SearchDomain(pattern.domain)
		if !reflect.DeepEqual(pattern.args, args) {
			t.Errorf("\n[%d]: expected reflect args: %v, got %v", i, pattern.args, args)
		}
		if err != pattern.err {
			t.Errorf("\n[%d]: expected error: %v, got %v", i, pattern.err, err)
		}
	}
}
