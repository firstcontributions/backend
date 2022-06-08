package service

// import (
// 	"context"
// 	"log"
// 	"net"

// 	"github.com/firstcontributions/backend/internal/grpc/users/proto"
// 	"github.com/firstcontributions/backend/internal/models/usersstore"
// 	"github.com/firstcontributions/backend/internal/models/usersstore/mongo"
// 	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/exporters/jaeger"
// 	"go.opentelemetry.io/otel/propagation"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	tracesdk "go.opentelemetry.io/otel/sdk/trace"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
// 	"google.golang.org/grpc"
// )

// type Service struct {
// 	*Config
// 	usersstore.Store
// 	*proto.UnimplementedUsersServiceServer
// }

// func (s *Service) Init(ctx context.Context) error {
// 	err := s.Config.DecodeEnv()
// 	if err != nil {
// 		return err
// 	}
// 	s.Store, err = mongo.NewUsersStore(ctx, *s.MongoURL)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func ListenAndServe(ctx context.Context) {
// 	s := &Service{
// 		Config: &Config{},
// 	}
// 	if err := s.Init(ctx); err != nil {
// 		log.Fatal(err)
// 	}
// 	lis, err := net.Listen("tcp", ":"+*s.Port)
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	tp, err := tracerProvider()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Register our TracerProvider as the global so any imported
// 	// instrumentation in the future will default to using it.
// 	otel.SetTracerProvider(tp)
// 	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

// 	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))

// 	// Register services
// 	proto.RegisterUsersServiceServer(grpcServer, s)

// 	log.Printf("GRPC server listening on %v", lis.Addr())

// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }

// func tracerProvider() (*tracesdk.TracerProvider, error) {
// 	// Create the Jaeger exporter
// 	exp, err := jaeger.New(jaeger.WithAgentEndpoint())
// 	if err != nil {
// 		return nil, err
// 	}
// 	tp := tracesdk.NewTracerProvider(
// 		// Always be sure to batch in production.
// 		tracesdk.WithBatcher(exp),
// 		// Record information about this application in a Resource.
// 		tracesdk.WithResource(resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceNameKey.String("userservice"),
// 			attribute.String("environment", "prod"),
// 			attribute.Int64("ID", 1),
// 		)),
// 	)
// 	return tp, nil
// }
