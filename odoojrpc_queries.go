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
	"strings"
)

// Common Odoo Queries
func (o *Odoo) ModelMap(model string, field string) (map[string]int, error) {
	ids := map[string]int{}
	rr, err := o.SearchRead(strings.Replace(model, "_", ".", -1), FilterArg{}, 0, 0, []string{field})
	if err != nil {
		return ids, err
	}
	for _, r := range rr {
		switch k := r[field].(type) {
		case string:
			switch v := r["id"].(type) {
			case float64:
				ids[k] = int(v)
			}
		}
	}
	return ids, nil
}

// CompanyID record
func (o *Odoo) CompanyID(companyName string) (int, error) {
	return o.GetID("res.company", []any{FilterArg{"name", "=", companyName}})
}

// PartnerID record
func (o *Odoo) PartnerID(partnerName string) (int, error) {
	return o.GetID("res.partner", []any{FilterArg{"name", "=", partnerName}})
}

// CountryID record
func (o *Odoo) CountryID(countryName string) (int, error) {
	return o.GetID("res.country", []any{FilterArg{"name", "=", countryName}})
}

// StateID record
func (o *Odoo) StateID(countryID int, stateName string) (int, error) {
	return o.GetID("res.country.state", []any{FilterArg{"name", "=", stateName}, FilterArg{"country_id", "=", countryID}})
}

// FiscalPosition record
func (o *Odoo) FiscalPosition(countryID int, fiscalName string) (int, error) {
	return o.GetID("account.fiscal.position", []any{FilterArg{"country_id", "=", countryID}, FilterArg{"name", "like", fiscalName}})
}
