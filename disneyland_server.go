package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/skygrid/disneyland/disneyland"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type DisneylandServerConfig struct {
	ServerCert  string `yaml:"server_cert"`
	ServerKey   string `yaml:"server_key"`
	CACert      string `yaml:"ca_cert"`
	ListenOn    string `yaml:"listen_on"`
	DatabaseURI string `yaml:"db_uri"`
}

const maxMessageSizeInBytes = 5 * 1024 * 1024

var Config *DisneylandServerConfig

func getTransportCredentials() (*credentials.TransportCredentials, error) {
	peerCert, err := tls.LoadX509KeyPair(Config.ServerCert, Config.ServerKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(Config.CACert)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tc := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{peerCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	return &tc, nil
}

func main() {
	Config = &DisneylandServerConfig{}
	config_path := os.Getenv("DISNEYLAND_CONFIG")
	content, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = yaml.Unmarshal([]byte(content), Config)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	storage, err := disneyland.NewDisneylandStorage(Config.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", Config.ListenOn)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := &disneyland.Server{
		Storage: storage,
	}

	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	logrusEntry := logrus.NewEntry(logger)
	transportCredentials, err := getTransportCredentials()
	if err != nil {
		log.Fatalf("failed to get credentials: %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMessageSizeInBytes),
		grpc.MaxSendMsgSize(maxMessageSizeInBytes),
		grpc.Creds(*transportCredentials),
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
	disneyland.RegisterDisneylandServer(s, server)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Print("Server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
