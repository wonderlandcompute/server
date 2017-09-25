package main

import (
	"github.com/sirupsen/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"gitlab.com/lambda-hse/optimus/optimus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"encoding/base64"
)

const (
	port = ":50051"
)

func main() {
	key_base64 := "o/nSETEx"
	key, err := base64.StdEncoding.DecodeString(key_base64)
	if err != nil {
		log.Fatal(err)
	}

	db_uri := os.Getenv("OPTIMUS_TEST_DB")
	storage, err := optimus.NewOptimusStorage(db_uri)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := &optimus.Server{
		Storage:      storage,
		SecretKey:    key,
	}
	server.Init()

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	logrusEntry := logrus.NewEntry(logger)

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_auth.UnaryServerInterceptor(nil),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			grpc_auth.StreamServerInterceptor(nil),
		),
	)
	optimus.RegisterOptimusServer(s, server)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Print("Server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
