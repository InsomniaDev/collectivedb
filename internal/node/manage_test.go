package node

import (
	"reflect"
	"testing"

	"github.com/insomniadev/collectivedb/internal/types"
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
		{
			name: "DataNew",
			args: args{
				update: &types.DataUpdate{
					DataUpdate: types.CollectiveDataUpdate{
						Update:     true,
						UpdateType: types.NEW,
						UpdateData: types.Data{
							ReplicaNodeGroup: 1,
							DataKey:          "12",
							Database:         "test",
						},
					},
				},
			},
		},
		{
			name: "DataUpdate",
			args: args{
				update: &types.DataUpdate{
					DataUpdate: types.CollectiveDataUpdate{
						Update:     true,
						UpdateType: types.UPDATE,
						UpdateData: types.Data{
							ReplicaNodeGroup: 1,
							DataKey:          "12",
							Database:         "testit",
						},
					},
				},
			},
		},
		{
			name: "DataDelete",
			args: args{
				update: &types.DataUpdate{
					DataUpdate: types.CollectiveDataUpdate{
						Update:     true,
						UpdateType: types.DELETE,
						UpdateData: types.Data{
							ReplicaNodeGroup: 1,
							DataKey:          "12",
							Database:         "testit",
						},
					},
				},
			},
		},
		{
			name: "ReplicaNew",
			args: args{
				update: &types.DataUpdate{
					ReplicaUpdate: types.CollectiveReplicaUpdate{
						Update:     true,
						UpdateType: types.NEW,
						UpdateReplica: types.ReplicaGroup{
							ReplicaNodeGroup:   1,
							SecondaryNodeGroup: 2,
							ReplicaNodes: []types.Node{
								{
									NodeId:    "1234",
									IpAddress: "123",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ReplicaUpdate",
			args: args{
				update: &types.DataUpdate{
					ReplicaUpdate: types.CollectiveReplicaUpdate{
						Update:     true,
						UpdateType: types.UPDATE,
						UpdateReplica: types.ReplicaGroup{
							ReplicaNodeGroup:   1,
							SecondaryNodeGroup: 2,
							ReplicaNodes: []types.Node{
								{
									NodeId:    "12345",
									IpAddress: "123",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ReplicaDelete",
			args: args{
				update: &types.DataUpdate{
					ReplicaUpdate: types.CollectiveReplicaUpdate{
						Update:     true,
						UpdateType: types.DELETE,
						UpdateReplica: types.ReplicaGroup{
							ReplicaNodeGroup:   1,
							SecondaryNodeGroup: 2,
							ReplicaNodes: []types.Node{
								{
									NodeId:    "1234",
									IpAddress: "123",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CollectiveUpdate(tt.args.update)

			switch tt.name {
			case "DataNew":
				if !reflect.DeepEqual(tt.args.update.DataUpdate.UpdateData, Collective.Data.DataLocations[0]) {
					t.Errorf("New Data = %v, want stored data to be = %v", tt.args.update.DataUpdate.UpdateData, Collective.Data.DataLocations[0])
				}
			case "DataUpdate":
				if !reflect.DeepEqual(tt.args.update.DataUpdate.UpdateData, Collective.Data.DataLocations[0]) {
					t.Errorf("Update Data = %v, want stored data to be = %v", tt.args.update.DataUpdate.UpdateData, Collective.Data.DataLocations[0])
				}
			case "DataDelete":
				if len(Collective.Data.DataLocations) != 0 {
					t.Errorf("Delete Data, expected 0 length, got %d", len(Collective.Data.DataLocations))
				}
			case "ReplicaNew":
				if !reflect.DeepEqual(tt.args.update.ReplicaUpdate.UpdateReplica, Collective.Data.CollectiveNodes[0]) {
					t.Errorf("New Data = %v, want stored data to be = %v", tt.args.update.ReplicaUpdate.UpdateReplica, Collective.Data.CollectiveNodes[0])
				}
			case "ReplicaUpdate":
				if !reflect.DeepEqual(tt.args.update.ReplicaUpdate.UpdateReplica, Collective.Data.CollectiveNodes[0]) {
					t.Errorf("Update Data = %v, want stored data to be = %v", tt.args.update.ReplicaUpdate.UpdateReplica, Collective.Data.CollectiveNodes[0])
				}
			case "ReplicaDelete":
				if len(Collective.Data.CollectiveNodes) != 0 {
					t.Errorf("Delete Data, expected 0 length, got %d", len(Collective.Data.CollectiveNodes))
				}
			}
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

func TestRetrieveSecondaryNodeGroupForDataEntry(t *testing.T) {
	Collective.Data.CollectiveNodes = []types.ReplicaGroup{
		{
			ReplicaNodeGroup:   1,
			SecondaryNodeGroup: 2,
		},
		{
			ReplicaNodeGroup: 2,
		},
	}

	type args struct {
		replicaNodeGroup *int
	}
	tests := []struct {
		name                   string
		args                   args
		wantSecondaryNodeGroup int
	}{
		// TODO: Add test cases.
		{
			name: "Has a node group",
			args: args{
				replicaNodeGroup: &Collective.Data.CollectiveNodes[0].ReplicaNodeGroup, // 1
			},
			wantSecondaryNodeGroup: 2,
		},
		{
			name: "Doesn't have a second node group",
			args: args{
				replicaNodeGroup: &Collective.Data.CollectiveNodes[1].ReplicaNodeGroup, // 2
			},
			wantSecondaryNodeGroup: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSecondaryNodeGroup := RetrieveSecondaryNodeGroupForDataEntry(tt.args.replicaNodeGroup); gotSecondaryNodeGroup != tt.wantSecondaryNodeGroup {
				t.Errorf("RetrieveSecondaryNodeGroupForDataEntry() = %v, want %v", gotSecondaryNodeGroup, tt.wantSecondaryNodeGroup)
			}
		})
	}
}
