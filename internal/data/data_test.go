package data

import (
	"reflect"
	"testing"
	"time"

	"github.com/insomniadev/collectivedb/internal/node"
	"github.com/insomniadev/collectivedb/internal/types"
)

func TestStoreDataInDatabase(t *testing.T) {

	key := "test"
	bucket := "test"
	data := []byte("hello")

	emptyKey := ""

	type args struct {
		key                *string
		bucket             *string
		data               *[]byte
		replicaStore       bool
		secondaryNodeGroup int
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *string
	}{
		{
			name: "Store without key",
			args: args{
				key:                &emptyKey,
				bucket:             &bucket,
				data:               &data,
				replicaStore:       false,
				secondaryNodeGroup: 0,
			},
			want:  true,
			want1: &key,
		},
		{
			name: "Stored successfully",
			args: args{
				key:                &key,
				bucket:             &bucket,
				data:               &data,
				replicaStore:       false,
				secondaryNodeGroup: 0,
			},
			want:  true,
			want1: &key,
		},
		{
			name: "Stored replica successfully",
			args: args{
				key:                &key,
				bucket:             &bucket,
				data:               &data,
				replicaStore:       true,
				secondaryNodeGroup: 0,
			},
			want:  true,
			want1: &key,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StoreDataInDatabase(tt.args.key, tt.args.bucket, tt.args.data, tt.args.replicaStore, tt.args.secondaryNodeGroup)
			if got != tt.want {
				t.Errorf("storeDataInDatabase() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				if tt.name == "Store without key" {
					if len(*got1) != len("1110b39e-14fd-4a20-b6ed-199709b14eac") {
						t.Errorf("storeDataInDatabase() got1 = %v, want %v", got1, tt.want1)
					}
				} else {
					t.Errorf("storeDataInDatabase() got1 = %v, want %v", got1, tt.want1)
				}
			}
		})
	}
}

func TestRetrieveDataFromDatabase(t *testing.T) {
	key := "test"
	bucket := "test"
	data := []byte("hello")
	StoreDataInDatabase(&key, &bucket, &data, false, 0)

	wrongKey := "don'texist"

	type args struct {
		key    *string
		bucket *string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *[]byte
	}{
		{
			name: "Retrieve Successfully",
			args: args{
				key:    &key,
				bucket: &bucket,
			},
			want:  true,
			want1: &data,
		},
		{
			name: "Doesnt exist",
			args: args{
				key:    &wrongKey,
				bucket: &bucket,
			},
			want:  false,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RetrieveDataFromDatabase(tt.args.key, tt.args.bucket)
			if got != tt.want {
				t.Errorf("RetrieveDataFromDatabase() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RetrieveDataFromDatabase() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeleteDataFromDatabase(t *testing.T) {
	key := "test"
	bucket := "test"
	data := []byte("hello")
	StoreDataInDatabase(&key, &bucket, &data, false, 0)

	wrongKey := "don'texist"

	type args struct {
		key    *string
		bucket *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Deleted",
			args: args{
				key:    &key,
				bucket: &bucket,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Doesn't exist",
			args: args{
				key:    &wrongKey,
				bucket: &bucket,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteDataFromDatabase(tt.args.key, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteDataFromDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("deleteDataFromDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetrieveAllReplicaData(t *testing.T) {
	key := "test"
	bucket := "test"
	data := []byte("hello")
	StoreDataInDatabase(&key, &bucket, &data, false, 0)
	count := 0

	node.Collective.Data.DataLocations = []types.Data{
		{
			ReplicaNodeGroup: 1,
			DataKey:          "test",
			Database:         "test",
		},
	}

	node.Collective.ReplicaNodeGroup = 1

	type args struct {
		inputData chan *types.StoredData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test successfully sending data",
			args: args{
				inputData: make(chan *types.StoredData),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func(inputData <-chan *types.StoredData) {
				for {
					data := <-inputData
					if data != nil {
						count++
					} else {
						return
					}
				}
			}(tt.args.inputData)

			RetrieveAllReplicaData(tt.args.inputData)
			time.Sleep(10 * time.Millisecond)
			if count != 1 {
				t.Errorf("Did not capture all data")
			}
		})
	}
}

func Test_distributeData(t *testing.T) {
	node.Collective = types.Controller{}
	node.Collective.ReplicaNodeGroup = 1
	node.Collective.Data.CollectiveNodes = []types.ReplicaGroup{
		{
			ReplicaNodeGroup: 1,
			ReplicaNodes: []types.Node{
				{
					NodeId:    "1",
					IpAddress: "127.0.0.1:9091",
				},
			},
		},
		{
			ReplicaNodeGroup: 2,
			ReplicaNodes: []types.Node{
				{
					NodeId:    "2",
					IpAddress: "127.0.0.1:9091",
				},
			},
		},
	}

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
			name: "Success_thisAsSecondaryNodeGroup",
			args: args{
				key:                &testKey,
				bucket:             &bucket,
				data:               &testValue,
				secondaryNodeGroup: 1,
			},
			wantErr: false,
		},
		{
			name: "Success_withSecondaryNodeGroup",
			args: args{
				key:                &testKey,
				bucket:             &bucket,
				data:               &testValue,
				secondaryNodeGroup: 2,
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
			if err := DistributeData(tt.args.key, tt.args.bucket, tt.args.data, tt.args.secondaryNodeGroup); (err != nil) != tt.wantErr {
				t.Errorf("distributeData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
