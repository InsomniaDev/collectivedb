package server

import (
	"github.com/insomniadev/collective-db/api/proto"
	"github.com/insomniadev/collective-db/internal/control"
)

func convertDataUpdatesToControlDataUpdate(incomingData *proto.DataUpdates) (convertedData *control.DataUpdate) {
	replicaNodes := []control.Node{}
	for i := range incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodes {
		replicaNodes = append(replicaNodes, control.Node{
			NodeId:    incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodes[i].NodeId,
			IpAddress: incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodes[i].IpAddress,
		})
	}

	return &control.DataUpdate{
		DataUpdate: control.CollectiveDataUpdate{
			Update:     incomingData.CollectiveUpdate.Update,
			UpdateType: int(incomingData.CollectiveUpdate.UpdateType),
			UpdateData: control.Data{
				ReplicaNodeGroup:  int(incomingData.CollectiveUpdate.Data.ReplicaNodeGroup),
				DataKey:           incomingData.CollectiveUpdate.Data.DataKey,
				Database:          incomingData.CollectiveUpdate.Data.Database,
				ReplicatedNodeIds: incomingData.CollectiveUpdate.Data.ReplicatedNodeIds,
			},
		},
		ReplicaUpdate: control.CollectiveReplicaUpdate{
			Update:     incomingData.ReplicaUpdate.Update,
			UpdateType: int(incomingData.ReplicaUpdate.UpdateType),
			UpdateReplica: control.ReplicaGroup{
				ReplicaNodeGroup: int(incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup),
				ReplicaNodes:     replicaNodes,
				FullGroup:        incomingData.ReplicaUpdate.UpdateReplica.FullGroup,
			},
		},
	}
}
