package main

import (
	"database/sql/driver"
)

// used in convert.go 88
func (c CountryType) Value() (driver.Value, error) {
	return c.String(), nil
}

// Goto convert.go 125
// Goto convert.go 241
func (c *CountryType) Scan(v interface{}) error {
	val, ok := v.([]byte)
	if !ok {
		panic("invalid postgres format for CountryType")
	}

	*c = pgCountryTypeToEnum[string(val)]
	return nil
}

var pgCountryTypeToEnum = map[string]CountryType{
	"us": CountryUS,
	"sg": CountrySG,
}
