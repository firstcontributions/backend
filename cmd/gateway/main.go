package main

import (
	"github.com/firstcontributions/backend/internal/gateway"
)

func main() {
	s := gateway.NewServer()
	if err := s.Init(); err != nil {
		panic(err)
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
