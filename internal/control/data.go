package control

import (
	"log"
	"os"

	"github.com/google/uuid"
	database "github.com/insomniadev/collective-db/internal/database"
)

// storeDataInDatabase
// will store the provided data into the database after checking if it requires an update first
// if the data belongs with a different replica group, it will send the update request to that replica group
func storeDataInDatabase(key, bucket *string, data *[]byte, replicaStore bool) (bool, *string) {
	ackLevel := os.Getenv("COLLECTIVE_ACK_LEVEL")
	if ackLevel == "" {
		ackLevel = "NONE"
	}

	var updateAndDistribute = func() bool {
		updated, key := database.Update(key, bucket, data)
		if !replicaStore {
			switch ackLevel {
			case "ALL":
				distributeData(key, bucket, data)
			case "NONE":
				go distributeData(key, bucket, data)
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

	// If the data doesn't exist yet, but a key was provided OR data exists and needs to be updated
	if dataVolume.DataKey == "" || dataVolume.ReplicaNodeGroup == controller.ReplicaNodeGroup {
		// Distribute the data across the collective
		return updateAndDistribute(), key
	}

	// Update the data on the different node
	for i := range controller.Data.CollectiveNodes {
		if controller.Data.CollectiveNodes[i].ReplicaNodeGroup == dataVolume.ReplicaNodeGroup {
			// TODO: API - Send the data to the leader for that replica group
			log.Println(controller.Data.CollectiveNodes[i].ReplicaNodes[0].IpAddress)
			// return the boolean from this call
			return false, nil
		}
	}
	return false, nil
}

// retrieveDataFromDatabase
func retrieveDataFromDatabase(key, bucket *string) (bool, *[]byte) {
	if exists, value := database.Get(key, bucket); exists {
		return exists, value
	}

	// The data does not exist on this node
	// TODO: Determine what node the data exists on
	// TODO: Go retrieve the data and then return it here

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
	// TODO: Determine what node the data exists on
	// TODO: Send a delete command to all replicas to delete the data
	// TODO: Wait for response of confirmation if that env variable is set

	return false, nil
}