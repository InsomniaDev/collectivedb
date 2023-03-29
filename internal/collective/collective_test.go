package collective

import (
	"log"
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/insomniadev/collective-db/internal/data"
	"github.com/insomniadev/collective-db/internal/node"
	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/proto/server"
	"github.com/insomniadev/collective-db/internal/types"
	"google.golang.org/grpc"
)

func init() {
	// Server type for working with the gRPC server
	// type grpcServerType struct {
	// 	proto.UnimplementedRouteGuideServer

	// 	dictionary_mu sync.Mutex
	// }

	// // Create and return the gRPC server
	// NewGrpcServer := func() *grpcServerType {
	// 	s := &grpcServerType{}
	// 	return s
	// }

	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterRouteGuideServer(grpcServer, server.NewGrpcServer())
	go grpcServer.Serve(lis)

	// Initial setup
	node.Active = true
	node.Collective.IpAddress = "127.0.0.1:9090"
	node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress = "127.0.0.1:9090"
}

func Test_createUuid(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Created",
			want: "9136b94f-552e-42d8-a7bc-0c5b8acf50df",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createUuid(); len(got) != len(tt.want) {
				t.Errorf("createUuid() = len(%d), want len(%d)", len(got), len(tt.want))
			}
		})
	}
}

func Test_determineIpAddress(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Url provided",
			want: "192.168.1.1",
		},
		{
			name: "Environment",
			want: "192-168-1-1.default.pod.cluster.local:9090",
		},
		{
			name: "Kubernetes",
			want: "192-168-1-1.default.pod.cluster.local:9090",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Url provided" {
				os.Setenv("COLLECTIVE_HOST_URL", "192.168.1.1")
			} else if tt.name == "Environment" {
				os.Setenv("COLLECTIVE_HOST_URL", "")
				os.Setenv("COLLECTIVE_IP", "192.168.1.1")
				os.Setenv("COLLECTIVE_RESOLVER_FILE", "test.conf")
			} else {
				os.Setenv("COLLECTIVE_IP", "192.168.1.1")
				os.Setenv("COLLECTIVE_RESOLVER_FILE", "test.conf")
			}
			if got := determineIpAddress(); got != tt.want {
				t.Errorf("determineIpAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_determineReplicas(t *testing.T) {

	node.Collective.Data.CollectiveNodes = []types.ReplicaGroup{
		{
			ReplicaNodeGroup: 1,
			ReplicaNodes: []types.Node{
				{
					NodeId:    "1",
					IpAddress: "127.0.0.1:9090",
				},
				{
					NodeId:    "2",
					IpAddress: "127.0.0.1:9090",
				},
				{
					NodeId:    "3",
					IpAddress: "127.0.0.1:9090",
				},
			},
			FullGroup: true,
		},
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Success1",
			wantErr: false,
		},
		{
			// Specified replica count is lower than controller amount
			name:    "Success2",
			wantErr: false,
		},
		{
			// Too high of a replica count to number of collective nodes
			name:    "Replicas",
			wantErr: true,
		},
		{
			// Failed to parse the number of replicas
			name:    "Failed",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Success2":
				os.Setenv("COLLECTIVE_REPLICA_COUNT", "2")
			case "Replicas":
				os.Setenv("COLLECTIVE_REPLICA_COUNT", "2")
				node.Collective.Data.CollectiveNodes[0].FullGroup = false
			case "Failed":
				os.Setenv("COLLECTIVE_REPLICA_COUNT", "NAN")
			}

			if err := determineReplicas(); (err != nil) != tt.wantErr {
				t.Errorf("determineReplicas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addToDataDictionary(t *testing.T) {
	node.Collective.Data.DataLocations = []types.Data{
		{
			ReplicaNodeGroup: 1,
			DataKey:          "1",
			Database:         "test",
		},
	}

	type args struct {
		dataToInsert types.Data
	}
	tests := []struct {
		name           string
		args           args
		wantUpdateType int
		wantUpdated    bool
	}{
		{
			name: "Updated",
			args: args{
				dataToInsert: types.Data{
					ReplicaNodeGroup: 2,
					DataKey:          "1",
					Database:         "test",
				},
			},
			wantUpdateType: types.UPDATE,
			wantUpdated:    true,
		},
		{
			name: "New",
			args: args{
				dataToInsert: types.Data{
					ReplicaNodeGroup: 2,
					DataKey:          "2",
					Database:         "test",
				},
			},
			wantUpdateType: types.NEW,
			wantUpdated:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUpdateType := node.AddToDataDictionary(tt.args.dataToInsert)
			if gotUpdateType != tt.wantUpdateType {
				t.Errorf("addToDataDictionary() gotUpdateType = %v, want %v", gotUpdateType, tt.wantUpdateType)
			}
		})
	}
}

func Test_removeDataFromSecondaryNodeGroup(t *testing.T) {
	node.Collective.Data.DataLocations = []types.Data{
		{
			ReplicaNodeGroup: 1,
			DataKey:          "1",
			Database:         "test",
		},
		{
			ReplicaNodeGroup: 2,
			DataKey:          "2",
			Database:         "test",
		},
	}
	node.Collective.ReplicaNodeGroup = 1
	node.Collective.Data.CollectiveNodes = []types.ReplicaGroup{
		{
			ReplicaNodeGroup: 1,
			ReplicaNodes: []types.Node{
				{
					NodeId:    "1",
					IpAddress: "127.0.0.1:9090",
				},
			},
			SecondaryNodeGroup: 2,
		},
		{
			ReplicaNodeGroup: 2,
			ReplicaNodes: []types.Node{
				{
					NodeId:    "2",
					IpAddress: "127.0.0.1:9090",
				},
			},
		},
	}

	type args struct {
		secondaryGroup int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				secondaryGroup: 2,
			},
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				secondaryGroup: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeDataFromSecondaryNodeGroup(tt.args.secondaryGroup); (err != nil) != tt.wantErr {
				t.Errorf("removeDataFromSecondaryNodeGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_terminateReplicas(t *testing.T) {
	node.Collective.ReplicaNodeGroup = 1
	key := "test"
	bucket := "test"
	newData := []byte("hello")
	replicaCount = 3
	data.StoreDataInDatabase(&key, &bucket, &newData, false, 0)
	node.Collective.ReplicaNodeIds = []string{node.Collective.NodeId}
	node.Collective.Data.DataLocations = []types.Data{
		{
			ReplicaNodeGroup: 1,
			DataKey:          "test",
			Database:         "test",
		},
	}
	node.Collective.Data.CollectiveNodes = []types.ReplicaGroup{
		{
			ReplicaNodeGroup: 1,
			ReplicaNodes: []types.Node{
				{
					NodeId:    "1",
					IpAddress: "127.0.0.1:9090",
				},
			},
			SecondaryNodeGroup: 2,
		},
		{
			ReplicaNodeGroup: 2,
			ReplicaNodes: []types.Node{
				{
					NodeId:    "2",
					IpAddress: "127.0.0.1:9090",
				},
			},
		},
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Last Node in Collective",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := terminateReplicas(); (err != nil) != tt.wantErr {
				t.Errorf("terminateReplicas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_retrieveDataDictionary(t *testing.T) {
	os.Setenv("COLLECTIVE_MAIN_BROKERS", "")

	node.Collective = types.Controller{}
	node.Collective.IpAddress = "127.0.0.1:9090"
	newCluster := []types.ReplicaGroup{
		{
			ReplicaNodeGroup:   1,
			SecondaryNodeGroup: 0,
			ReplicaNodes: []types.Node{
				{
					NodeId:    node.Collective.NodeId,
					IpAddress: "127.0.0.1:9090",
				},
			},
			FullGroup: false,
		},
	}

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "New_Cluster",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node.Collective.IpAddress = "127.0.0.1:9090"
			retrieveDataDictionary()
			switch tt.name {
			case "New_Cluster":
				if !reflect.DeepEqual(node.Collective.Data.CollectiveNodes, newCluster) {
					t.Errorf("retrieveDataDictionary() got = %v, want %v", node.Collective.Data.CollectiveNodes, newCluster)
				}
			}
		})
	}
}
