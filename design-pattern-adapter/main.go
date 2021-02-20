package main

import (
	"fmt"
	"time"
)

type ConditionAdapter interface {
	Validate() bool
}

func main() {
	a := []ConditionAdapter{}
	time, _ := time.Parse("2006-01-02 15:04:05 MST", "2021-02-20 18:11:00"+" WIB")
	datetimeCondition := &DatetimeCondition{
		Operator: "greater_than",
		Value:    time,
	}
	weekdayCondition := &WeekdayCondition{
		Operator: "contains",
		Value:    []int64{1, 2, 6},
	}
	a = append(a, datetimeCondition)
	a = append(a, weekdayCondition)

	result := false
	allRes := []bool{}
	for i, v := range a {
		value := v.Validate()
		if i == 0 {
			result = value
		} else {
			result = result && value
		}
		allRes = append(allRes, value)
	}

	fmt.Println("result", result)
	fmt.Println("all res", allRes)
}
