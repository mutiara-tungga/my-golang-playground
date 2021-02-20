package main

import "time"

type WeekdayCondition struct {
	Operator string
	Value    []int64
}

// Validate to validate if now.Weekdays is contains on WeekdayCondition.Value
func (w *WeekdayCondition) Validate() bool {
	now := time.Now()
	weekday := int64(now.Weekday())

	for _, v := range w.Value {
		if weekday == v {
			return true
		}
	}

	return false
}
