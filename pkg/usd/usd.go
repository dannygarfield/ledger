package usd

import (
	"fmt"
	"strconv"
	"strings"
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
	if strings.Contains(s, ".") {
		float, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return USD(-1), err
		}
		int := int(float*100)
		return USD(int), nil
	} else {
		return USD(-1), fmt.Errorf("Could not parse %s, must enter dollar amount with decimal and cents", s)
	}
}
