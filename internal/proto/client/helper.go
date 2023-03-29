package client

import (
	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/types"
)

func convertDataUpdatesToControlDataUpdate(incomingData *proto.DataUpdates) (convertedData *types.DataUpdate) {
	replicaNodes := []types.Node{}
	for i := range incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodes {
		replicaNodes = append(replicaNodes, types.Node{
			NodeId:    incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodes[i].NodeId,
			IpAddress: incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodes[i].IpAddress,
		})
	}

	return &types.DataUpdate{
		DataUpdate: types.CollectiveDataUpdate{
			Update:     incomingData.CollectiveUpdate.Update,
			UpdateType: int(incomingData.CollectiveUpdate.UpdateType),
			UpdateData: types.Data{
				ReplicaNodeGroup: int(incomingData.CollectiveUpdate.Data.ReplicaNodeGroup),
				DataKey:          incomingData.CollectiveUpdate.Data.DataKey,
				Database:         incomingData.CollectiveUpdate.Data.Database,
			},
		},
		ReplicaUpdate: types.CollectiveReplicaUpdate{
			Update:     incomingData.ReplicaUpdate.Update,
			UpdateType: int(incomingData.ReplicaUpdate.UpdateType),
			UpdateReplica: types.ReplicaGroup{
				ReplicaNodeGroup:   int(incomingData.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup),
				ReplicaNodes:       replicaNodes,
				FullGroup:          incomingData.ReplicaUpdate.UpdateReplica.FullGroup,
				SecondaryNodeGroup: int(incomingData.ReplicaUpdate.UpdateReplica.SecondaryNodeGroup),
			},
		},
	}
}