package main

import (
	"context"

	"github.com/firstcontributions/backend/internal/grpc/users/service"
)

func main() {
	service.ListenAndServe(context.Background())
}
