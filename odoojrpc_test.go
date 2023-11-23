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
	"fmt"
	"testing"
)

var urlPatterns = []struct {
	schema        string
	hostname      string
	port          int
	expected      string
	expectedError error
}{
	{"", "", 0, "", fmt.Errorf("init error: %w", ErrSchema)},
	{"", "localhost", 0, "", fmt.Errorf("init error: %w", ErrSchema)},
	{"", "", 8069, "", fmt.Errorf("init error: %w", ErrSchema)},
	{"http", "", 0, "", fmt.Errorf("init error: %w", ErrPort)},
	{"http", "localhost", 0, "", fmt.Errorf("init error: %w", ErrPort)},
	{"http", "", 8069, "", fmt.Errorf("init error: %w", ErrHostLen)},
	{"https", "", 0, "", fmt.Errorf("init error: %w", ErrPort)},
	{"https", "localhost", 0, "", fmt.Errorf("init error: %w", ErrPort)},
	{"https", "", 8069, "", fmt.Errorf("init error: %w", ErrHostLen)},
	{"ftp", "", 0, "", fmt.Errorf("init error: %w", ErrSchema)},
	{"ftp", "localhost", 0, "", fmt.Errorf("init error: %w", ErrSchema)},
	{"ftp", "", 8069, "", fmt.Errorf("init error: %w", ErrSchema)},
	{"ftp", "localhost", 8069, "", fmt.Errorf("init error: %w", ErrSchema)},
	{"http", "localhost", 8069, "http://localhost:8069/jsonrpc", nil},
	{"https", "localhost", 8069, "https://localhost:8069/jsonrpc", nil},
}

func TestURL(t *testing.T) {
	for i, pattern := range urlPatterns {
		o := Odoo{
			Hostname: pattern.hostname,
			Port:     pattern.port,
			Schema:   pattern.schema,
		}

		o.genURL()

		if len(o.URL) != len(pattern.expected) {
			t.Errorf("\n[%d]: slice size not equal, expected: %d, got %d", i, len(pattern.expected), len(o.URL))
			t.Errorf("\n[%d]: expected %s, got %s", i, pattern.expected, o.URL)
		}
	}
}

func TestInitError(t *testing.T) {
	for i, pattern := range urlPatterns {
		o := new(Odoo)
		o.Hostname = pattern.hostname
		o.Port = pattern.port
		o.Schema = pattern.schema
		err := o.Init()
		if err != nil {
			if err.Error() != pattern.expectedError.Error() {
				t.Errorf("\n[%d]: expected %s, got %s", i, pattern.expectedError, err)
			}
		}
	}
}
