package odoojrpc

import (
	"fmt"
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
	{"('name','=','My Name')", FilterArg{"name", "=", "My Name"}, nil},
	{"('name','like','My Name')", FilterArg{"name", "like", "My Name"}, nil},
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
	for _, pattern := range searchDomainPatterns {
		fmt.Println("test: domain:", pattern.domain, "pattern.err:", pattern.err)
		args, err := SearchDomain(pattern.domain)
		if !reflect.DeepEqual(pattern.args, args) {
			t.Errorf("\nexpected reflect args: %v, got %v", pattern.args, args)
		} else {
			fmt.Println("pass: args are equal")
		}
		if err != pattern.err {
			t.Errorf("\nexpected error: %v, got %v", pattern.err, err)
		} else {
			fmt.Println("pass: errors are equal")
		}
		fmt.Println()
	}
}
