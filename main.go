package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/insomniadev/collective-db/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 50051, "The server port")
)

// Setup the api
// https://github.com/grpc/grpc-go/blob/master/examples/route_guide/server/server.go
type server struct {
	pb.UnimplementedRouteGuideServer
}

func newServer() *server {
	s := &server{}
	return s
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = "x509/server_cert.pem"
		}
		if *keyFile == "" {
			*keyFile = "x509/server_key.pem"
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterRouteGuideServer(grpcServer, newServer())
	grpcServer.Serve(lis)

	// Retrieve the depth that the hash function should extend to
	// depth := os.Getenv("DEPTH")
	// num, err := strconv.Atoi(depth)
	// if err != nil {
	// 	// Non numeric environment variable input, setting to 1
	// 	fmt.Println("Depth incorrectly input, setting to 1")
	// 	num = 1
	// }

	// fmt.Println(database.FNV32a("test"))

	// db := database.Init(num)
	// fmt.Println(db.TotalDepth)
}
