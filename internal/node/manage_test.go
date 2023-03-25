package node

import (
	"reflect"
	"testing"

	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/types"
)

func Test_removeFromDictionarySlice(t *testing.T) {
	type args struct {
		s []types.ReplicaGroup
		i int
	}
	tests := []struct {
		name string
		args args
		want []types.ReplicaGroup
	}{
		{
			name: "Remove an element",
			args: args{
				s: []types.ReplicaGroup{
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
			want: []types.ReplicaGroup{
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
				s: []types.ReplicaGroup{
					{
						ReplicaNodeGroup: 1,
					},
				},
				i: 0,
			},
			want: []types.ReplicaGroup{},
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

func TestCollectiveUpdate(t *testing.T) {
	type args struct {
		update *types.DataUpdate
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CollectiveUpdate(tt.args.update)
		})
	}
}

func TestRetrieveFromDataDictionary(t *testing.T) {
	key := "1"
	doesntExistKey := "2"
	Collective.Data.DataLocations = []types.Data{
		{
			ReplicaNodeGroup: 1,
			DataKey:          key,
			Database:         "test",
		},
	}

	type args struct {
		key *string
	}
	tests := []struct {
		name     string
		args     args
		wantData types.Data
	}{
		{
			name: "Exists",
			args: args{
				key: &key,
			},
			wantData: Collective.Data.DataLocations[0],
		},
		{
			name: "Doesn't Exist",
			args: args{
				key: &doesntExistKey,
			},
			wantData: types.Data{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotData := RetrieveFromDataDictionary(tt.args.key); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("retrieveFromDataDictionary() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func Test_sendClientUpdateDictionaryRequest(t *testing.T) {
	type args struct {
		ipAddress *string
		update    *proto.DataUpdates
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
			if err := SendClientUpdateDictionaryRequest(tt.args.ipAddress, tt.args.update); (err != nil) != tt.wantErr {
				t.Errorf("sendClientUpdateDictionaryRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
