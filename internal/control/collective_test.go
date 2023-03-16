package control

import (
	"os"
	"reflect"
	"testing"
)

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
			want: "192-168-1-1.default.pod.cluster.local",
		},
		{
			name: "Kubernetes",
			want: "192-168-1-1.default.pod.cluster.local",
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

	controller.Data.CollectiveNodes = []ReplicaGroup{
		{
			ReplicaNodeGroup: 1,
			ReplicaNodes: []Node{
				{
					NodeId:    "1",
					IpAddress: "1",
				},
				{
					NodeId:    "2",
					IpAddress: "2",
				},
				{
					NodeId:    "3",
					IpAddress: "3",
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
				controller.Data.CollectiveNodes[0].FullGroup = false
			case "Failed":
				os.Setenv("COLLECTIVE_REPLICA_COUNT", "NAN")
			}

			if err := determineReplicas(); (err != nil) != tt.wantErr {
				t.Errorf("determineReplicas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_distributeData(t *testing.T) {
	// TODO: Add in tests for the secondaryNodeGroup

	bucket := "test"

	testKey := "key"
	testValue := []byte("value")

	falseBucket := ""
	falseKey := ""

	type args struct {
		key                *string
		bucket             *string
		data               *[]byte
		secondaryNodeGroup int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				key:                &testKey,
				bucket:             &bucket,
				data:               &testValue,
				secondaryNodeGroup: 0,
			},
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				key:                &falseKey,
				bucket:             &falseBucket,
				data:               &testValue,
				secondaryNodeGroup: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := distributeData(tt.args.key, tt.args.bucket, tt.args.data, tt.args.secondaryNodeGroup); (err != nil) != tt.wantErr {
				t.Errorf("distributeData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_addToDataDictionary(t *testing.T) {
	controller.Data.DataLocations = []Data{
		{
			ReplicaNodeGroup: 1,
			DataKey:          "1",
			Database:         "test",
			ReplicatedNodeIds: []string{
				"1", "2", "3", "5",
			},
		},
	}

	type args struct {
		dataToInsert Data
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
				dataToInsert: Data{
					ReplicaNodeGroup: 2,
					DataKey:          "1",
					Database:         "test",
				},
			},
			wantUpdateType: UPDATE,
			wantUpdated:    true,
		},
		{
			name: "New",
			args: args{
				dataToInsert: Data{
					ReplicaNodeGroup: 2,
					DataKey:          "2",
					Database:         "test",
				},
			},
			wantUpdateType: NEW,
			wantUpdated:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUpdateType := addToDataDictionary(tt.args.dataToInsert)
			if gotUpdateType != tt.wantUpdateType {
				t.Errorf("addToDataDictionary() gotUpdateType = %v, want %v", gotUpdateType, tt.wantUpdateType)
			}
		})
	}
}

func Test_retrieveFromDataDictionary(t *testing.T) {
	key := "1"
	doesntExistKey := "2"
	controller.Data.DataLocations = []Data{
		{
			ReplicaNodeGroup: 1,
			DataKey:          key,
			Database:         "test",
			ReplicatedNodeIds: []string{
				"1", "2", "3", "5",
			},
		},
	}

	type args struct {
		key *string
	}
	tests := []struct {
		name     string
		args     args
		wantData Data
	}{
		{
			name: "Exists",
			args: args{
				key: &key,
			},
			wantData: controller.Data.DataLocations[0],
		},
		{
			name: "Doesn't Exist",
			args: args{
				key: &doesntExistKey,
			},
			wantData: Data{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotData := retrieveFromDataDictionary(tt.args.key); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("retrieveFromDataDictionary() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func Test_removeFromDictionarySlice(t *testing.T) {
	type args struct {
		s []ReplicaGroup
		i int
	}
	tests := []struct {
		name string
		args args
		want []ReplicaGroup
	}{
		{
			name: "Remove an element",
			args: args{
				s: []ReplicaGroup{
					{
						ReplicaNodeGroup: 1,
					},
					{
						ReplicaNodeGroup: 2,
					},
					{
						ReplicaNodeGroup: 3,
					},
				},
				i: 1,
			},
			want: []ReplicaGroup{
				{
					ReplicaNodeGroup: 1,
				},
				{
					ReplicaNodeGroup: 3,
				},
			},
		},
		{
			name: "An empty array",
			args: args{
				s: []ReplicaGroup{
					{
						ReplicaNodeGroup: 1,
					},
				},
				i: 0,
			},
			want: []ReplicaGroup{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeFromDictionarySlice(tt.args.s, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeFromDictionarySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDataFromSecondaryNodeGroup(t *testing.T) {
	type args struct {
		secondaryGroup int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeDataFromSecondaryNodeGroup(tt.args.secondaryGroup); (err != nil) != tt.wantErr {
				t.Errorf("removeDataFromSecondaryNodeGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
