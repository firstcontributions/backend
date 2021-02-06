package main

import (
	"fmt"
	"log"

	"net"

	"github.com/firstcontributions/backend/internal/profile"
	"github.com/firstcontributions/backend/internal/profile/proto"
	"github.com/firstcontributions/backend/pkg/userctx"

	"google.golang.org/grpc"
)

func main() {

	s := profile.NewService()
	if err := s.Init(); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *s.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(userctx.Authorize),
	)

	proto.RegisterProfileServiceServer(grpcServer, s)

	log.Printf("listening at :%s", *s.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
