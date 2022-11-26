package odoojrpc

// Common Odoo Queries

// CompanyID record
func (o *Odoo) CompanyID(companyName string) int {
	return o.GetID("res.company", []any{FilterArg{"name", "=", companyName}})
}

// PartnerID record
func (o *Odoo) PartnerID(partnerName string) int {
	return o.GetID("res.partner", []any{FilterArg{"name", "=", partnerName}})
}

// CountryID record
func (o *Odoo) CountryID(countryName string) int {
	return o.GetID("res.country", []any{FilterArg{"name", "=", countryName}})
}

// StateID record
func (o *Odoo) StateID(countryID int, stateName string) int {
	return o.GetID("res.country.state", []any{FilterArg{"name", "=", stateName}, FilterArg{"country_id", "=", countryID}})
}

// FiscalPosition record
func (o *Odoo) FiscalPosition(countryID int, fiscalName string) int {
	return o.GetID("account.fiscal.position", []any{FilterArg{"country_id", "=", countryID}, FilterArg{"name", "like", fiscalName}})
}
