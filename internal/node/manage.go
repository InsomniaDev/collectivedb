package node

import (
	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/proto/client"
	"github.com/insomniadev/collective-db/internal/types"
)

// CollectiveUpdate
//
// This will go through and update the collective memory, not touching the actual data
func CollectiveUpdate(update *types.DataUpdate) {
	CollectiveMemoryMutex.Lock()

	// If this is a data update
	if update.DataUpdate.Update {
		// Update the data dictionary
		switch update.DataUpdate.UpdateType {
		case types.NEW:
			// Adds the new element to the end of the array
			Collective.Data.DataLocations = append(Collective.Data.DataLocations, update.DataUpdate.UpdateData)
		case types.UPDATE:
			// Updates the element where it is
			for i := range Collective.Data.DataLocations {
				if Collective.Data.DataLocations[i].DataKey == update.DataUpdate.UpdateData.DataKey {
					Collective.Data.DataLocations[i] = update.DataUpdate.UpdateData
					break
				}
			}
		case types.DELETE:
			// Deletes the element from the array
			for i := range Collective.Data.DataLocations {
				if Collective.Data.DataLocations[i].DataKey == update.DataUpdate.UpdateData.DataKey {
					Collective.Data.DataLocations = removeFromDictionarySlice(Collective.Data.DataLocations, i)
					break
				}
			}
		}
	} else if update.ReplicaUpdate.Update {
		// Update the data dictionary
		switch update.ReplicaUpdate.UpdateType {
		case types.NEW:
			// Adds the new element to the end of the array
			Collective.Data.CollectiveNodes = append(Collective.Data.CollectiveNodes, update.ReplicaUpdate.UpdateReplica)
		case types.UPDATE:
			// Updates the element where it is
			for i := range Collective.Data.CollectiveNodes {
				if Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup {
					Collective.Data.CollectiveNodes[i] = update.ReplicaUpdate.UpdateReplica
					break
				}
			}
		case types.DELETE:
			// Deletes the element from the array
			for i := range Collective.Data.CollectiveNodes {
				if Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup {
					Collective.Data.CollectiveNodes = removeFromDictionarySlice(Collective.Data.CollectiveNodes, i)
					break
				}
			}
		}
	}

	CollectiveMemoryMutex.Unlock()
}

// removeFromDictionarySlice
// removes the specified index from the slice and returns that slice, this does reorder the array by switching out the elements
func removeFromDictionarySlice[T types.Collective](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// DictionaryUpdate
//
//	Will update this node with the incoming collective information from the other nodes
func DictionaryUpdate(update *types.DataUpdate) {
	// Update this node with the incoming information
	// Send the data to the first url in the next replica group
	// Send the new data to all replicas in this replica group

	// If we get data here, then this is the leader of the replica group (or the first node in the array)
	CollectiveUpdate(update)

	// Create an array of ReplicaNodes info
	replicaNodes := []*proto.ReplicaNodes{}
	for i := range update.ReplicaUpdate.UpdateReplica.ReplicaNodes {
		replicaNodes = append(replicaNodes, &proto.ReplicaNodes{
			NodeId:    update.ReplicaUpdate.UpdateReplica.ReplicaNodes[i].NodeId,
			IpAddress: update.ReplicaUpdate.UpdateReplica.ReplicaNodes[i].IpAddress,
		})
	}

	// Assemble the data to be sent
	protoData := &proto.DataUpdates{
		CollectiveUpdate: &proto.CollectiveDataUpdate{
			Update:     update.DataUpdate.Update,
			UpdateType: int32(update.DataUpdate.UpdateType),
			Data: &proto.CollectiveData{
				ReplicaNodeGroup: int32(update.DataUpdate.UpdateData.ReplicaNodeGroup),
				DataKey:          update.DataUpdate.UpdateData.DataKey,
				Database:         update.DataUpdate.UpdateData.Database,
			},
		},
		ReplicaUpdate: &proto.CollectiveReplicaUpdate{
			Update:     update.ReplicaUpdate.Update,
			UpdateType: int32(update.ReplicaUpdate.UpdateType),
			UpdateReplica: &proto.UpdateReplica{
				ReplicaNodeGroup: int32(update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup),
				FullGroup:        update.ReplicaUpdate.UpdateReplica.FullGroup,
				ReplicaNodes:     replicaNodes,
			},
		},
	}

	// Send the data onward
	// Update the replicas in this replica group
	for i := range Collective.ReplicaNodes {
		// Let's not send to ourselves here
		if Collective.ReplicaNodes[i].NodeId != Collective.NodeId {
			// Send the data to the replica node

			// Send the data to the api endpoint for the ReplicaUpdate function - ReplicaUpdate rpc
			// initialize first call
			updateDictionary := make(chan *proto.DataUpdates)
			client.ReplicaUpdate(&Collective.ReplicaNodes[i].IpAddress, updateDictionary)

			// Send the data into the dictionary update function
			updateDictionary <- protoData
			updateDictionary <- nil
		}
	}

	// Send to the next replica group in the list
	for i := range Collective.Data.CollectiveNodes {
		// Go to where this node group is in the array
		if Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == Collective.ReplicaNodeGroup {
			// Send to the next replica group in the list
			// 	Check that there is another element in the array
			//  Confirm that we aren't going to send to the replicanodegroup that started this request
			if len(Collective.Data.CollectiveNodes) >= i+2 && Active &&
				(Collective.Data.CollectiveNodes[i].ReplicaNodeGroup != update.DataUpdate.UpdateData.ReplicaNodeGroup ||
					Collective.Data.CollectiveNodes[i].ReplicaNodeGroup != update.ReplicaUpdate.UpdateReplica.ReplicaNodeGroup) {

				// initialize first call
				updateDictionary := make(chan *proto.DataUpdates)
				client.DictionaryUpdate(&Collective.Data.CollectiveNodes[i+1].ReplicaNodes[0].IpAddress, updateDictionary)

				// Send the data into the dictionary update function
				updateDictionary <- protoData
				updateDictionary <- nil

				// Break out of the loop and allow the next process to send the data, otherwise all data will always be sent from one location
				break
			}
		}
	}
}

// AddToDataDictionary
//
//	Will add the data structure to the dictionary, or update the location
func AddToDataDictionary(dataToInsert types.Data) (updateType int) {
	CollectiveMemoryMutex.Lock()

	for i := range Collective.Data.DataLocations {
		if Collective.Data.DataLocations[i].DataKey == dataToInsert.DataKey {
			// already exists, so check if the data matches
			Collective.Data.DataLocations[i] = dataToInsert

			// Unlock and return
			CollectiveMemoryMutex.Unlock()
			return types.UPDATE
		}
	}

	// if the data doesn't exist already
	Collective.Data.DataLocations = append(Collective.Data.DataLocations, dataToInsert)

	// Unlock and return
	CollectiveMemoryMutex.Unlock()
	return types.NEW
}

// RetrieveFromDataDictionary
// Will retrieve the key from the dictionary if it exists
func RetrieveFromDataDictionary(key *string) (data types.Data) {

	for i := range Collective.Data.DataLocations {
		if Collective.Data.DataLocations[i].DataKey == *key {
			return Collective.Data.DataLocations[i]
		}
	}

	return
}

// SendClientUpdateDictionaryRequest
//
// extracted function that is used to send the update without all of the additional boilerplate code everywhere
func SendClientUpdateDictionaryRequest(ipAddress *string, update *proto.DataUpdates) error {

	// Create the channel
	updateDictionary := make(chan *proto.DataUpdates)

	// Call the dictionary function before passing the data into the channel
	// send the update to the first node in the list
	go client.DictionaryUpdate(ipAddress, updateDictionary)

	updateDictionary <- update
	updateDictionary <- nil
	return nil
}
