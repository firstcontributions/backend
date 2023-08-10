package main

import (
	"fmt"

	"github.com/firstcontributions/backend/internal/gateway"
)

func main() {
	fmt.Println("from gateway new CI pipeline")
	s := gateway.NewServer()
	if err := s.Init(); err != nil {
		panic(err)
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
