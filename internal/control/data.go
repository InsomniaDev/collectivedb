package control

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/insomniadev/collective-db/api/client"
	"github.com/insomniadev/collective-db/api/proto"
	database "github.com/insomniadev/collective-db/internal/database"
)

// storeDataInDatabase
// will store the provided data into the database after checking if it requires an update first
// if the data belongs with a different replica group, it will send the update request to that replica group
func storeDataInDatabase(key, bucket *string, data *[]byte, replicaStore bool, secondaryNodeGroup int) (bool, *string) {
	ackLevel := os.Getenv("COLLECTIVE_ACK_LEVEL")
	if ackLevel == "" {
		ackLevel = "NONE"
	}

	var updateAndDistribute = func() bool {
		updated, key := database.Update(key, bucket, data)
		if !replicaStore {
			switch ackLevel {
			case "ALL":
				distributeData(key, bucket, data, secondaryNodeGroup)
			case "NONE":
				go distributeData(key, bucket, data, secondaryNodeGroup)
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
	dataVolume := retrieveFromDataDictionary(key)

	// If the data doesn't exist yet, but a key was provided OR data exists and needs to be updated OR this data was sent in with a secondaryNodeGroup equal to this one
	if dataVolume.DataKey == "" || dataVolume.ReplicaNodeGroup == controller.ReplicaNodeGroup || controller.ReplicaNodeGroup == secondaryNodeGroup {
		// Distribute the data across the collective
		return updateAndDistribute(), key
	}

	// Update the data on the different node
	for i := range controller.Data.CollectiveNodes {
		if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == dataVolume.ReplicaNodeGroup {
			// Send the data to the leader for that replica group - DataUpdate rpc
			log.Println(controller.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress)

			protoData := make(chan *proto.Data)
			err := client.ReplicaDataUpdate(&controller.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress, protoData)

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

// retrieveAllReplicaData
// Will return with all of the replicated data
func retrieveAllReplicaData(inputData chan<- *StoredData) {
	for i := range controller.Data.DataLocations {
		if controller.Data.DataLocations[i].ReplicaNodeGroup == controller.ReplicaNodeGroup {
			if exists, data := retrieveDataFromDatabase(&controller.Data.DataLocations[i].DataKey, &controller.Data.DataLocations[i].Database); exists {
				inputData <- &StoredData{
					ReplicaNodeGroup: controller.ReplicaNodeGroup,
					DataKey:          controller.Data.DataLocations[i].DataKey,
					Database:         controller.Data.DataLocations[i].Database,
					Data:             *data,
				}
			}
		}
	}
	inputData <- nil
}

// retrieveDataFromDatabase
func retrieveDataFromDatabase(key, bucket *string) (bool, *[]byte) {
	if exists, value := database.Get(key, bucket); exists {
		return exists, value
	}

	// The data does not exist on this node
	// Determine what node the data exists on
	for i := range controller.Data.DataLocations {
		if controller.Data.DataLocations[i].DataKey == *key {

			// Go retrieve the data and then return it here - GetData rpc
			for j := range controller.Data.CollectiveNodes {
				if controller.Data.CollectiveNodes[j].ReplicaNodeGroup == controller.Data.DataLocations[i].ReplicaNodeGroup {

					data, err := client.GetData(&controller.Data.CollectiveNodes[j].ReplicaNodes[0].IpAddress, &proto.Data{
						Key:              *key,
						Database:         *bucket,
						ReplicaNodeGroup: int32(controller.Data.DataLocations[i].ReplicaNodeGroup),
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

func deleteDataFromDatabase(key, bucket *string) (bool, error) {
	if deleted, err := database.Delete(key, bucket); deleted {
		return true, err
	} else if err != nil {
		return false, err
	}

	// The data does not exist on this node
	// Determine what node the data exists on
	for i := range controller.Data.DataLocations {
		if controller.Data.DataLocations[i].DataKey == *key {

			// Send delete command to the first node in the replica group that contains the data
			deleteData := make(chan *proto.Data)
			if err := client.DeleteData(&controller.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress, deleteData); err != nil {
				return false, err
			}

			deleteData <- &proto.Data{
				Key:              *key,
				Database:         *bucket,
				ReplicaNodeGroup: int32(controller.Data.DataLocations[i].ReplicaNodeGroup),
			}
			deleteData <- nil
			return true, nil
		}
	}

	return false, nil
}
