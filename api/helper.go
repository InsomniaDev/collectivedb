package api

import "github.com/insomniadev/collective-db/internal/control"

func convertDataUpdatesToControlDataUpdate(incomingData *DataUpdates) (convertedData *control.DataUpdate) {
	return &control.DataUpdate{
		DataUpdate:    {},
		ReplicaUpdate: {},
	}
}
