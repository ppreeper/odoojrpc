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
	"reflect"
	"testing"
)

var searchDomainPatterns = []struct {
	domain string
	args   []any
	err    error
}{
	{"", []any{}, nil},
	{"('')", []any{}, errSyntax},
	{"('','')", []any{}, errSyntax},
	{"('a','=')", []any{}, errSyntax},
	{"('name')", []any{}, errSyntax},
	{"('name','=')", []any{}, errSyntax},
	{"('name','=','My Name')", []any{[]any{"name", "=", "My Name"}}, nil},
	{"('name','like','My Name')", []any{[]any{"name", "like", "My Name"}}, nil},
	{"('name','=','My Name'),('name','=','My Name')", []any{}, errSyntax},
	{"[('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}}, nil},
	{"[('name','=','My Name'),('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}}, nil},
	{"[('name','=','My Name'),'!',('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"!", []any{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),'&',('name','=','My Name'),('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"&", []any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"|", []any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'!',('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}, []any{"!", []any{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'&',('name','=','My Name')]", []any{}, errSyntax},
	{"[('name','=','My Name'),('name','=','My Name'),'&',('name','=','My Name'),('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}, []any{"&", []any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name')]", []any{}, errSyntax},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}, []any{"|", []any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name'),'!',('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}, []any{"|", []any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}}, []any{"!", []any{"name", "=", "My Name"}}}, nil},
	{"[('name','=','My Name'),('name','=','My Name'),'|',('name','=','My Name'),('name','=','My Name'),'!',('name','=','My Name'),('name','=','My Name')]", []any{[]any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}, []any{"|", []any{"name", "=", "My Name"}, []any{"name", "=", "My Name"}}, []any{"!", []any{"name", "=", "My Name"}}, []any{"name", "=", "My Name"}}, nil},
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
