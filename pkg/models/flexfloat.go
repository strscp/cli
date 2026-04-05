package models

import (
	"encoding/json"
	"strconv"
)

// FlexFloat is a float64 that can unmarshal from both JSON numbers and strings.
// Laravel's decimal cast returns strings like "4.25" to preserve precision.
type FlexFloat float64

// UnmarshalJSON handles both "3.14" and 3.14.
func (f *FlexFloat) UnmarshalJSON(data []byte) error {
	// Try number first
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexFloat(num)
		return nil
	}

	// Try string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = FlexFloat(num)
	return nil
}

// Float64 returns the underlying float64 value.
func (f FlexFloat) Float64() float64 {
	return float64(f)
}
