package main

import (
	"time"
)

type DatetimeCondition struct {
	Operator string
	Value    time.Time
}

// Validate is to validate if now is equal or greather than or less than DatetimeCondition.Value
func (d *DatetimeCondition) Validate() bool {
	now := time.Now()

	switch d.Operator {
	case "equal":
		return now.Equal(d.Value)
	case "greater_than":
		return now.After(d.Value)
	case "less_than":
		return now.Before(d.Value)
	}

	return false
}
