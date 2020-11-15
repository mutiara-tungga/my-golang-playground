package main

import (
	"context"
	"fmt"
)

type ctxKey string

// CobaContext coba
type CobaContext struct {
	context context.Context
}

// FuncCobaContext coba
func FuncCobaContext() {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, ctxKey("X-Request-ID"), "hahahaha")
	fmt.Println(ctx.Value(ctxKey("X-Request-ID")).(string))

	a := CobaContext{
		// context: ctx,
	}
	fmt.Println(a.context == nil)
	// reqID := a.context.Value(ctxKey("X-Request-ID")).(string)
	// fmt.Println(reqID)
}

func main() {
	FuncCobaContext()
}
