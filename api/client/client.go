package client

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/insomniadev/collective-db/api/proto"
	"github.com/insomniadev/collective-db/resources"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tls = ""
)

func init() {
	tls = os.Getenv("COLLECTIVE_TLS")
}

// getConnectionOptions
// will create the option config for the node
func getConnectionOptions(ipAddress *string) *[]grpc.DialOption {
	var opts []grpc.DialOption
	if tls != "" {
		caFile := resources.Path("x509/ca_cert.pem")
		creds, err := credentials.NewClientTLSFromFile(caFile, "")
		if err != nil {
			log.Fatalf("Failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return &opts
}

// SyncDataRequest
// Is responsible for syncing all data from the other node in the replica group
func SyncDataRequest(ipAddress *string, data chan<- *proto.Data) error {
	// Setup the client
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if syncClient, err := client.SyncDataRequest(ctx, &proto.SyncIp{IpAddress: *ipAddress}); err != nil {
		for {
			in, err := syncClient.Recv()
			if err == io.EOF {
				// read done.
				return nil
			}
			if err != nil {
				return err
			}
			data <- in
		}
	}

	return nil
}

// DataUpdate
// Will attempt to update the data from the specified location, which should be the replication node group leader
func DataUpdate(ipAddress *string, dataChan <-chan *proto.Data) error {
	// Setup the client
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.DataUpdate(ctx)

	if err != nil {
		log.Fatalf("stream.RecordRoute failed: %v", err)
	}

	go func() {
		for {
			data := <-dataChan
			if data == nil {
				if err := stream.CloseSend(); err != nil {
					log.Println(err)
				}
				break
			}

			if err := stream.Send(data); err != nil {
				log.Printf("stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			}
		}
	}()

	return nil
}

// GetData
// Will attempt to retrieve the data from the specified location, which should be the replication node group leader
func GetData(ipAddress *string, data *proto.Data) (*proto.Data, error) {
	// Setup the client
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to retrieve the data from the endpoint
	returnedData, err := client.GetData(ctx, data)
	if returnedData != nil || err != nil {
		return nil, err
	} else {
		return returnedData, nil
	}
}

// DeleteData
// Will take an array of data fields and have them deleted from the provided ipaddress
func DeleteData(ipAddress *string, dataChan <-chan *proto.Data) error {
	// Setup the client
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.DeleteData(ctx)

	if err != nil {
		log.Fatalf("stream.RecordRoute failed: %v", err)
	}

	go func() {
		for {
			data := <-dataChan
			if data == nil {
				if err := stream.CloseSend(); err != nil {
					log.Println(err)
				}
				break
			}

			if err := stream.Send(data); err != nil {
				log.Printf("stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
				break
			}
		}
	}()
	return nil
}

// DictionaryUpdate
// Will update the dictionary with the provided dataset
func DictionaryUpdate(ipAddress *string, dataChan <-chan *proto.DataUpdates) error {
	// Setup the client
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.DictionaryUpdate(ctx)

	if err != nil {
		log.Printf("stream.RecordRoute failed: %v", err)
		return err
	}

	go func() {
		for {
			data := <-dataChan
			if data == nil {
				if err := stream.CloseSend(); err != nil {
					log.Println(err)
				}
				break
			}

			if err := stream.Send(data); err != nil {
				log.Printf("stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
				break
			}
		}
	}()

	return nil
}

// ReplicaUpdate
// Will have a collective update to the attached replica through the replica update point
func ReplicaUpdate(ipAddress *string, dataChan <-chan *proto.DataUpdates) error {
	// Setup the client
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ReplicaUpdate(ctx)

	if err != nil {
		log.Fatalf("stream.RecordRoute failed: %v", err)
	}

	go func() {
		for {
			data := <-dataChan
			if data == nil {
				if err := stream.CloseSend(); err != nil {
					log.Println(err)
				}
				break
			}

			if err := stream.Send(data); err != nil {
				log.Printf("stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			}
		}
	}()

	return nil
}

// ReplicaDataUpdate
// Will send the data to be stored on the connected replicas
func ReplicaDataUpdate(ipAddress *string, dataChan <-chan *proto.Data) error {
	// Setup the client
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ReplicaDataUpdate(ctx)

	if err != nil {
		log.Fatalf("stream.RecordRoute failed: %v", err)
	}

	go func() {
		for {
			data := <-dataChan
			if data == nil {
				if err := stream.CloseSend(); err != nil {
					log.Println(err)
				}
				break
			}

			if err := stream.Send(data); err != nil {
				log.Printf("stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			}
		}
	}()

	return nil
}
