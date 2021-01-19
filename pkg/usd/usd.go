package usd

import (
	"fmt"
	"strconv"
)

// USD represents a quantity of US money.
//
// It is stored as an integer quantity of cents.
type USD int

// String formats this USD in the conventional way.
// String implements the flag.Value interface
func (d *USD) String() string {
	if d == nil {
		return fmt.Sprint(nil)
	}
	// format &USD(12304) as $123.04
	return fmt.Sprintf("$%d.%.2d", *d/100, *d%100)
}

func StringToUsd(s string) (USD, error) {
	int, err := strconv.Atoi(s)
	if err != nil {
		return USD(-1), err
	}
	return USD(int), nil
}
