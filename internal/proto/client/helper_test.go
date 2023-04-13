package client

import (
	"reflect"
	"testing"

	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/types"
)

func Test_convertDataUpdatesToControlDataUpdate(t *testing.T) {
	type args struct {
		incomingData *proto.DataUpdates
	}
	tests := []struct {
		name              string
		args              args
		wantConvertedData *types.DataUpdate
	}{
		{
			name: "Successful Data",
			args: args{
				incomingData: &proto.DataUpdates{
					CollectiveUpdate: &proto.CollectiveDataUpdate{
						Update:     true,
						UpdateType: types.NEW,
						Data: &proto.CollectiveData{
							ReplicaNodeGroup: 1,
							DataKey:          "test",
							Database:         "test",
						},
					},
				},
			},
			wantConvertedData: &types.DataUpdate{
				DataUpdate: types.CollectiveDataUpdate{
					Update:     true,
					UpdateType: types.NEW,
					UpdateData: types.Data{
						ReplicaNodeGroup: 1,
						DataKey:          "test",
						Database:         "test",
					},
				},
			},
		},
		{
			name: "Successful Replica",
			args: args{
				incomingData: &proto.DataUpdates{
					ReplicaUpdate: &proto.CollectiveReplicaUpdate{
						Update:     true,
						UpdateType: types.NEW,
						UpdateReplica: &proto.UpdateReplica{
							ReplicaNodeGroup: 1,
							FullGroup:        false,
							ReplicaNodes: []*proto.ReplicaNodes{
								{
									NodeId:    "1",
									IpAddress: "127.0.0.1",
								},
							},
							SecondaryNodeGroup: 2,
						},
					},
				},
			},
			wantConvertedData: &types.DataUpdate{
				ReplicaUpdate: types.CollectiveReplicaUpdate{
					Update:     true,
					UpdateType: types.NEW,
					UpdateReplica: types.ReplicaGroup{
						ReplicaNodeGroup:   1,
						SecondaryNodeGroup: 2,
						ReplicaNodes: []types.Node{
							{
								NodeId:    "1",
								IpAddress: "127.0.0.1",
							},
						},
						FullGroup: false,
					},
				},
			},
		},
		{
			name: "Successful Both",
			args: args{
				incomingData: &proto.DataUpdates{
					CollectiveUpdate: &proto.CollectiveDataUpdate{
						Update:     true,
						UpdateType: types.NEW,
						Data: &proto.CollectiveData{
							ReplicaNodeGroup: 1,
							DataKey:          "test",
							Database:         "test",
						},
					},
					ReplicaUpdate: &proto.CollectiveReplicaUpdate{
						Update:     true,
						UpdateType: types.NEW,
						UpdateReplica: &proto.UpdateReplica{
							ReplicaNodeGroup: 1,
							FullGroup:        false,
							ReplicaNodes: []*proto.ReplicaNodes{
								{
									NodeId:    "1",
									IpAddress: "127.0.0.1",
								},
							},
							SecondaryNodeGroup: 2,
						},
					},
				},
			},
			wantConvertedData: &types.DataUpdate{
				DataUpdate: types.CollectiveDataUpdate{
					Update:     true,
					UpdateType: types.NEW,
					UpdateData: types.Data{
						ReplicaNodeGroup: 1,
						DataKey:          "test",
						Database:         "test",
					},
				},
				ReplicaUpdate: types.CollectiveReplicaUpdate{
					Update:     true,
					UpdateType: types.NEW,
					UpdateReplica: types.ReplicaGroup{
						ReplicaNodeGroup:   1,
						SecondaryNodeGroup: 2,
						ReplicaNodes: []types.Node{
							{
								NodeId:    "1",
								IpAddress: "127.0.0.1",
							},
						},
						FullGroup: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConvertedData := ConvertDataUpdatesToControlDataUpdate(tt.args.incomingData); !reflect.DeepEqual(gotConvertedData, tt.wantConvertedData) {
				t.Errorf("convertDataUpdatesToControlDataUpdate() = %v, want %v", gotConvertedData, tt.wantConvertedData)
			}
		})
	}
}
