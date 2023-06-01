package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/insomniadev/collectivedb/internal/proto"
	"github.com/insomniadev/collectivedb/internal/types"
	"github.com/insomniadev/collectivedb/resources"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tls = ""

	storeDataProm = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "client_data_update",
		Help: "Gauge of the data as it is updated",
	})

	retrieveDataProm = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "client_data_retrieval",
		Help: "Gauge of the data as it is retrieved",
	})

	deleteDataProm = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "client_data_deletion",
		Help: "Gauge of the data as it is deleted",
	})
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

// SyncCollectiveRequest
// Is responsible for syncing all collective data in the cluster to this new node
func SyncCollectiveRequest(ipAddress *string, data chan<- *types.DataUpdate) error {
	// Setup the client
	log.Println(fmt.Scanf("SyncCollectiveRequest:%s", ipAddress))
	connOpts := getConnectionOptions(ipAddress)
	conn, err := grpc.Dial(*ipAddress, *connOpts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewRouteGuideClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if syncClient, err := client.SyncCollectiveRequest(ctx, &proto.SyncIp{IpAddress: *ipAddress}); err != nil {
		for {
			in, err := syncClient.Recv()
			if err == io.EOF {
				// read done.
				data <- nil
				return nil
			}
			if err != nil {
				data <- nil
				return err
			}
			data <- ConvertDataUpdatesToControlDataUpdate(in)
		}
	}

	return nil
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
				data <- nil
				return nil
			}
			if err != nil {
				data <- nil
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
		log.Fatalf("DataUpdate stream.RecordRoute failed: %v", err)
	}

	for data := range dataChan {
		if err := stream.Send(data); err != nil {
			log.Printf("DataUpdate stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			break
		}
	}

	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
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
		log.Fatalf("DeleteData stream.RecordRoute failed: %v", err)
	}

	for data := range dataChan {
		if err := stream.Send(data); err != nil {
			log.Printf("DeleteData stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			break
		}
	}

	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
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
	stream, err := client.DictionaryUpdate(ctx)
	defer cancel()

	if err != nil {
		log.Printf("DictionaryUpdate stream.RecordRoute failed: %v", err)
		return err
	}

	for data := range dataChan {
		if err := stream.Send(data); err != nil {
			log.Printf("DictionaryUpdate stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			break
		}
	}

	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
	return nil
}

// ReplicaUpdate
// Will have a collective update to the attached replica through the replica update point
func ReplicaUpdate(ipAddress *string, dataChan <-chan *proto.DataUpdates) error {
	// TODO: Need to determine a way to notice a lost node and automatically trigger a data distribution

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
		log.Fatalf("ReplicaUpdate stream.RecordRoute failed: %v", err)
	}

	for data := range dataChan {
		if err := stream.Send(data); err != nil {
			log.Printf("ReplicaUpdate stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			break
		}
	}

	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
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
		log.Fatalf("ReplicaDataUpdate stream.RecordRoute failed: %v", err)
	}

	for data := range dataChan {
		if err := stream.Send(data); err != nil {
			log.Printf("ReplicaDataUpdate stream.RecordRoute: stream.Send(%v) failed: %v", data, err)
			break
		}
	}

	if err := stream.CloseSend(); err != nil {
		log.Println(err)
	}
	return nil
}
