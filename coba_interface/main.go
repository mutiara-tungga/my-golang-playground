package main

import "fmt"

type Haha1 interface {
	Test1() string
}

var _ Haha1 = (*Test)(nil)

type Haha2 interface {
	Test1() string
	Test2() string
}

var _ Haha2 = (*Test)(nil)

type ImplHaha1 struct {
	haha1 Haha1
}

func (h ImplHaha1) PrintHaha() {
	fmt.Println("ini adalah implementasi Haha 1 : ", h.haha1.Test1())
}

type ImplHaha2 struct {
	haha2 Haha2
}

func (h ImplHaha2) PrintHaha() {
	fmt.Println("ini adalah implementasi Haha 2 : ", h.haha2.Test1())
	fmt.Println("ini adalah implementasi Haha 2 : ", h.haha2.Test2())
}

func main() {
	ih1 := ImplHaha1{
		haha1: Test{},
	}

	ih1.PrintHaha()

	ih2 := ImplHaha2{
		haha2: Test{},
	}

	ih2.PrintHaha()
}
