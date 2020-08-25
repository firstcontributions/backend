package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"net"
	"os"

	"github.com/firstcontributions/firstcontributions/internal/profile"
	"github.com/firstcontributions/firstcontributions/internal/proto"

	"google.golang.org/grpc"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func main() {

	s := profile.NewService()
	if err := s.Init(); err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *s.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterProfileServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
