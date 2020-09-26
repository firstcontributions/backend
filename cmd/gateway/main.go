package main

import (
	"log"

	"github.com/firstcontributions/backend/internal/gateway"
)

func main() {
	s := gateway.NewServer()
	if err := s.Init(); err != nil {
		panic(err)
	}
	log.Print("-------------------------", s.BlockKey)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
