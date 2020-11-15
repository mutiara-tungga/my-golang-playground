package main

import (
	"fmt"
	"regexp"
)

func main() {
	a := "testing_123"

	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	fmt.Println(re.MatchString("123"))
	fmt.Println(re.MatchString("abc"))
	fmt.Println(re.MatchString(a))
	fmt.Println(re.MatchString("世界"))
}
