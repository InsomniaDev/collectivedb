package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/proto/server"
	"github.com/insomniadev/collective-db/resources"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 9090, "The port for data insertion; defaults 9090")
	nodePort = flag.Int("nodePort", 9091, "The port for collective communication; defaults 9091")
)

// Setup the api
// https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = resources.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = resources.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterRouteGuideServer(grpcServer, server.NewGrpcServer())
	grpcServer.Serve(lis)
}

// TODO: set a block to catch a SIG KILL or SIG TERM signal that will then fire off the terminateReplicas functionality
