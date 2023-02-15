package client

import (
	"context"
	"fmt"
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

// DeleteData
// Will take an array of data fields and have them deleted from the provided ipaddress
func DeleteData(ipAddress *string, data *proto.DataArray) error {
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

	deleted, err := client.DeleteData(ctx, data)
	if !deleted.UpdatedSuccessfully || err != nil {
		return err
	} else {
		return nil
	}
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
		log.Fatalf("stream.RecordRoute failed: %v", err)
	}

	for {
		data := <-dataChan
		if data == nil {
			if err := stream.CloseSend(); err != nil {
				return err
			}
			break
		}

		if err := stream.Send(data); err != nil {
			return fmt.Errorf("stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
		}
	}

	return nil
}
