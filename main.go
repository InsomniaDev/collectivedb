package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/insomniadev/collective-db/api"
	"github.com/insomniadev/collective-db/internal/collective"
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

	sigs       = make(chan os.Signal, 1)
	grpcServer *grpc.Server
)

// Setup the api
// https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go

func main() {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleAppClose()

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
	grpcServer = grpc.NewServer(opts...)
	proto.RegisterRouteGuideServer(grpcServer, server.NewGrpcServer())
	go grpcServer.Serve(lis)

	quitAfterTenSeconds := 0
	for {
		if collective.IsActive() {
			break
		} else {
			time.Sleep(1 * time.Second)
			quitAfterTenSeconds++
			if quitAfterTenSeconds > 10 {
				panic("Application failed to get healthy")
			}
			log.Println("initializing")
		}
	}
	api.Start()
}

// set a block to catch a SIG KILL or SIG TERM signal that will then fire off the terminateReplicas functionality
func handleAppClose() {
	sig := <-sigs
	log.Printf("Received signal (%s) to quit application\n", sig.String())
	if deactivated := collective.Deactivate(); deactivated {
		log.Println("Successfully shutdown node and redistributed the traffic")
	} else {
		log.Println("Failed to successfully shutdown the node")
	}
	grpcServer.Stop()
	api.Stop()
	os.Exit(0)
}
