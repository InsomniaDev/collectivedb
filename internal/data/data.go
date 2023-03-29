package data

import (
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
	database "github.com/insomniadev/collective-db/internal/database"
	"github.com/insomniadev/collective-db/internal/node"
	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/proto/client"
	"github.com/insomniadev/collective-db/internal/types"
)

func DistributeData(key, bucket *string, data *[]byte, secondaryNodeGroup int) error {

	if *key == "" || *bucket == "" {
		return errors.New("invalid parameters")
	}

	newData := types.Data{
		ReplicaNodeGroup: node.Collective.ReplicaNodeGroup,
		DataKey:          *key,
		Database:         *bucket,
	}

	if node.Active {
		// Create the data object to be sent
		dataUpdate := &proto.Data{
			Key:                *key,
			Database:           *bucket,
			Data:               *data,
			ReplicaNodeGroup:   int32(node.Collective.ReplicaNodeGroup),
			SecondaryNodeGroup: int32(secondaryNodeGroup),
		}

		// Send to each replica attached to this replica node group
		for i := range node.Collective.ReplicaNodes {
			if node.Collective.ReplicaNodes[i].NodeId != node.Collective.NodeId {
				updateReplica := make(chan *proto.Data)
				client.ReplicaDataUpdate(&node.Collective.ReplicaNodes[i].IpAddress, updateReplica)
				updateReplica <- dataUpdate
				updateReplica <- nil
			}
		}

		// Double check that the secondaryNodeGroup is 0 before starting to process
		if secondaryNodeGroup != 0 {
			for i := range node.Collective.Data.CollectiveNodes {
				if node.Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == secondaryNodeGroup {
					for j := range node.Collective.Data.CollectiveNodes[i].ReplicaNodes {
						updateReplica := make(chan *proto.Data)
						client.ReplicaDataUpdate(&node.Collective.Data.CollectiveNodes[i].ReplicaNodes[j].IpAddress, updateReplica)
						updateReplica <- dataUpdate
						updateReplica <- nil
					}

					// IF this replicaGroup is not complete and has a secondaryNodeGroup, THEN forward to all nodes in that group as well
					if !node.Collective.Data.CollectiveNodes[i].FullGroup {
						for j := range node.Collective.Data.CollectiveNodes {
							if node.Collective.Data.CollectiveNodes[j].ReplicaNodeGroup == node.Collective.Data.CollectiveNodes[i].SecondaryNodeGroup {
								// Send the update to the first node of that replica to start the update process from there
								dataUpdate.SecondaryNodeGroup = int32(node.Collective.Data.CollectiveNodes[i].SecondaryNodeGroup)

								updateReplica := make(chan *proto.Data)
								client.ReplicaDataUpdate(&node.Collective.Data.CollectiveNodes[j].ReplicaNodes[0].IpAddress, updateReplica)
								updateReplica <- dataUpdate
								updateReplica <- nil
								break
							}
						}
						break
					}
				}
			}
		}

		// Only update the data dictionary with this data if it was not sent here as part of the secondaryNodeGroup
		if secondaryNodeGroup != node.Collective.ReplicaNodeGroup {

			// Add this node to the DataDictionary
			updateType := node.AddToDataDictionary(newData)

			if err := node.SendClientUpdateDictionaryRequest(&node.Collective.Data.CollectiveNodes[0].ReplicaNodes[0].IpAddress, &proto.DataUpdates{
				CollectiveUpdate: &proto.CollectiveDataUpdate{
					Update:     true,
					UpdateType: int32(updateType),
					Data: &proto.CollectiveData{
						ReplicaNodeGroup: int32(newData.ReplicaNodeGroup),
						DataKey:          newData.DataKey,
						Database:         newData.Database,
					},
				},
			}); err != nil {
				return err
			}
		}

	}

	return nil
}

// StoreDataInDatabase
// will store the provided data into the database after checking if it requires an update first
// if the data belongs with a different replica group, it will send the update request to that replica group
func StoreDataInDatabase(key, bucket *string, data *[]byte, replicaStore bool, secondaryNodeGroup int) (bool, *string) {
	ackLevel := os.Getenv("COLLECTIVE_ACK_LEVEL")
	if ackLevel == "" {
		ackLevel = "NONE"
	}

	var updateAndDistribute = func() bool {
		updated, key := database.Update(key, bucket, data)
		if !replicaStore {
			switch ackLevel {
			case "ALL":
				DistributeData(key, bucket, data, secondaryNodeGroup)
			case "NONE":
				// FIXME: Need to set so we don't do race conditions, but also allow to return without confirmations from replicas... hmmm...
				DistributeData(key, bucket, data, secondaryNodeGroup)
			}
		}
		return updated
	}

	// Data is new and doesn't exist
	// Create a unique key and update since this is new data
	if *key == "" {
		newKey := uuid.New().String()
		key = &newKey

		// Distribute the data across the collective
		return updateAndDistribute(), key
	}

	// This data exists already
	// Determine what node the data is on, if the data does exist on a node
	dataVolume := node.RetrieveFromDataDictionary(key)

	// If the data doesn't exist yet, but a key was provided OR data exists and needs to be updated OR this data was sent in with a secondaryNodeGroup equal to this one
	if dataVolume.DataKey == "" || dataVolume.ReplicaNodeGroup == node.Collective.ReplicaNodeGroup || node.Collective.ReplicaNodeGroup == secondaryNodeGroup {
		// Distribute the data across the collective
		return updateAndDistribute(), key
	}

	// Update the data on the different node
	for i := range node.Collective.Data.CollectiveNodes {
		if node.Collective.Data.CollectiveNodes[i].ReplicaNodeGroup == dataVolume.ReplicaNodeGroup {
			// Send the data to the leader for that replica group - DataUpdate rpc
			log.Println(node.Collective.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress)

			protoData := make(chan *proto.Data)
			err := client.ReplicaDataUpdate(&node.Collective.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress, protoData)

			protoData <- &proto.Data{
				Key:              *key,
				Database:         *bucket,
				Data:             *data,
				ReplicaNodeGroup: int32(dataVolume.ReplicaNodeGroup),
			}
			protoData <- nil

			// return the boolean from this call
			if err != nil {
				return false, nil
			} else {
				return true, key
			}
		}
	}
	return false, nil
}

// RetrieveAllReplicaData
// Will return with all of the replicated data
func RetrieveAllReplicaData(inputData chan<- *types.StoredData) {
	for i := range node.Collective.Data.DataLocations {
		if node.Collective.Data.DataLocations[i].ReplicaNodeGroup == node.Collective.ReplicaNodeGroup {
			if exists, data := RetrieveDataFromDatabase(&node.Collective.Data.DataLocations[i].DataKey, &node.Collective.Data.DataLocations[i].Database); exists {
				inputData <- &types.StoredData{
					ReplicaNodeGroup: node.Collective.ReplicaNodeGroup,
					DataKey:          node.Collective.Data.DataLocations[i].DataKey,
					Database:         node.Collective.Data.DataLocations[i].Database,
					Data:             *data,
				}
			}
		}
	}
	inputData <- nil
}

// RetrieveDataFromDatabase
func RetrieveDataFromDatabase(key, bucket *string) (bool, *[]byte) {
	if exists, value := database.Get(key, bucket); exists {
		return exists, value
	}

	// The data does not exist on this node
	// Determine what node the data exists on
	for i := range node.Collective.Data.DataLocations {
		if node.Collective.Data.DataLocations[i].DataKey == *key {

			// Go retrieve the data and then return it here - GetData rpc
			for j := range node.Collective.Data.CollectiveNodes {
				if node.Collective.Data.CollectiveNodes[j].ReplicaNodeGroup == node.Collective.Data.DataLocations[i].ReplicaNodeGroup {

					data, err := client.GetData(&node.Collective.Data.CollectiveNodes[j].ReplicaNodes[0].IpAddress, &proto.Data{
						Key:              *key,
						Database:         *bucket,
						ReplicaNodeGroup: int32(node.Collective.Data.DataLocations[i].ReplicaNodeGroup),
					})
					if err != nil {
						return false, nil
					}
					return true, &data.Data
				}
			}
		}
	}

	// If data is not found on one replica node, it should attempt to pull from at least two nodes before declaring data doesn't exist
	// This should attempt to grab data from the ReplicatedNodes list

	return false, nil
}

func DeleteDataFromDatabase(key, bucket *string) (bool, error) {
	if deleted, err := database.Delete(key, bucket); deleted {
		return true, err
	} else if err != nil {
		return false, err
	}

	// The data does not exist on this node
	// Determine what node the data exists on
	for i := range node.Collective.Data.DataLocations {
		if node.Collective.Data.DataLocations[i].DataKey == *key {

			// Send delete command to the first node in the replica group that contains the data
			deleteData := make(chan *proto.Data)
			if err := client.DeleteData(&node.Collective.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress, deleteData); err != nil {
				return false, err
			}

			deleteData <- &proto.Data{
				Key:              *key,
				Database:         *bucket,
				ReplicaNodeGroup: int32(node.Collective.Data.DataLocations[i].ReplicaNodeGroup),
			}
			deleteData <- nil
			return true, nil
		}
	}

	return false, nil
}