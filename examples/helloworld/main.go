package main

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, main *app) error {
		fmt.Println("Hello, World!")
		return nil
	})
}
