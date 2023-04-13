package server

// TODO: Add unit tests for this file

import (
	context "context"
	"io"
	"log"
	"sync"

	"github.com/insomniadev/collective-db/internal/data"
	"github.com/insomniadev/collective-db/internal/database"
	"github.com/insomniadev/collective-db/internal/node"
	"github.com/insomniadev/collective-db/internal/proto"
	"github.com/insomniadev/collective-db/internal/proto/client"
	"github.com/insomniadev/collective-db/internal/types"
)

// Server type for working with the gRPC server
type grpcServer struct {
	// proto.UnimplementedRouteGuideServer
	proto.RouteGuideServer

	dictionary_mu sync.Mutex
}

// Create and return the gRPC server
func NewGrpcServer() *grpcServer {
	s := &grpcServer{}
	return s
}

// ReplicaUpdate receives a stream of data updates and
// returns with a boolean on if updated successfully
func (s *grpcServer) ReplicaUpdate(stream proto.RouteGuide_ReplicaUpdateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Lock the dictionary so that an update can occur here
		s.dictionary_mu.Lock()
		node.CollectiveUpdate(client.ConvertDataUpdatesToControlDataUpdate(in))
		s.dictionary_mu.Unlock()

		if err := stream.Send(&proto.Updated{
			UpdatedSuccessfully: true,
		}); err != nil {
			return err
		}
	}
}

// SyncCollectiveRequest
// Will send a request to the server to pull in all of the collective data to the newly joined node
func (s *grpcServer) SyncCollectiveRequest(syncIpAddress *proto.SyncIp, stream proto.RouteGuide_SyncCollectiveRequestServer) (err error) {
	// Make a channel that the process can yield the discovered data through
	storedData := make(chan *proto.DataUpdates)

	// Setup a process to wait for the returned data
	var wg sync.WaitGroup

	// Go through and return the data as it is discovered
	wg.Add(1)
	go func(chan *proto.DataUpdates) {
		defer wg.Done()
		for {
			data := <-storedData
			if data != nil {
				if err = stream.Send(data); err != nil {
					// Stop processing since we hit an error
					log.Println("SyncCollectiveRequest: ", err)
					return
				}
			} else {
				// We are done processing
				return
			}
		}
	}(storedData)

	// Cycle through the data and return
	node.CollectiveMemoryMutex.RLock()
	for i := range node.Collective.Data.CollectiveNodes {
		replicaNodes := []*proto.ReplicaNodes{}
		for j := range node.Collective.Data.CollectiveNodes[i].ReplicaNodes {
			replicaNodes = append(replicaNodes, &proto.ReplicaNodes{
				NodeId:    node.Collective.Data.CollectiveNodes[i].ReplicaNodes[j].NodeId,
				IpAddress: node.Collective.Data.CollectiveNodes[i].ReplicaNodes[j].IpAddress,
			})
		}

		storedData <- &proto.DataUpdates{
			ReplicaUpdate: &proto.CollectiveReplicaUpdate{
				Update:     true,
				UpdateType: types.NEW,
				UpdateReplica: &proto.UpdateReplica{
					ReplicaNodeGroup:   int32(node.Collective.Data.CollectiveNodes[i].ReplicaNodeGroup),
					SecondaryNodeGroup: int32(node.Collective.Data.CollectiveNodes[i].SecondaryNodeGroup),
					FullGroup:          node.Collective.Data.CollectiveNodes[i].FullGroup,
					ReplicaNodes:       replicaNodes,
				},
			},
		}
	}
	for i := range node.Collective.Data.DataLocations {
		storedData <- &proto.DataUpdates{
			CollectiveUpdate: &proto.CollectiveDataUpdate{
				Update:     true,
				UpdateType: types.NEW,
				Data: &proto.CollectiveData{
					ReplicaNodeGroup: int32(node.Collective.Data.DataLocations[i].ReplicaNodeGroup),
					DataKey:          node.Collective.Data.DataLocations[i].DataKey,
					Database:         node.Collective.Data.DataLocations[i].Database,
				},
			},
		}
	}
	node.CollectiveMemoryMutex.RUnlock()
	storedData <- nil

	// Wait until all stream data has been sent
	wg.Wait()
	return err
}

// SyncDataRequest
// Will send a request to the server to pull in all of the data to the newly joined node
func (s *grpcServer) SyncDataRequest(syncIpAddress *proto.SyncIp, stream proto.RouteGuide_SyncDataRequestServer) (err error) {
	// Make a channel that the process can yield the discovered data through
	storedData := make(chan *types.StoredData)

	// Setup a process to wait for the returned data
	var wg sync.WaitGroup

	// Go through and return the data as it is discovered
	wg.Add(1)
	go func(chan *types.StoredData) {
		defer wg.Done()
		for {
			data := <-storedData
			if data != nil {
				if err = stream.Send(&proto.Data{
					Key:              data.DataKey,
					Database:         data.Database,
					Data:             data.Data,
					ReplicaNodeGroup: int32(data.ReplicaNodeGroup),
				}); err != nil {
					// Stop processing since we hit an error
					return
				}
			} else {
				// We are done processing
				return
			}
		}
	}(storedData)

	// Request the data to be returned
	node.CollectiveMemoryMutex.RLock()
	for i := range node.Collective.Data.DataLocations {
		if node.Collective.Data.DataLocations[i].ReplicaNodeGroup == node.Collective.ReplicaNodeGroup {
			if exists, value := database.Get(&node.Collective.Data.DataLocations[i].DataKey, &node.Collective.Data.DataLocations[i].Database); exists {
				storedData <- &types.StoredData{
					ReplicaNodeGroup: node.Collective.ReplicaNodeGroup,
					DataKey:          node.Collective.Data.DataLocations[i].DataKey,
					Database:         node.Collective.Data.DataLocations[i].Database,
					Data:             *value,
				}
			}
		}
	}
	node.CollectiveMemoryMutex.RUnlock()
	storedData <- nil

	// Wait until all stream data has been sent
	wg.Wait()
	return err
}

// DictionaryUpdate
// Will send a stream of data entries that requie an update, will respond with a boolean for each entry sent
func (s *grpcServer) DictionaryUpdate(stream proto.RouteGuide_DictionaryUpdateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println(err)
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}

		// Lock the dictionary so that an update can occur here
		s.dictionary_mu.Lock()
		// control.CollectiveUpdate(convertDataUpdatesToControlDataUpdate(in))
		node.DictionaryUpdate(client.ConvertDataUpdatesToControlDataUpdate(in))
		s.dictionary_mu.Unlock()

		if err := stream.Send(&proto.Updated{
			UpdatedSuccessfully: true,
		}); err != nil {
			return err
		}
	}
}

// DataUpdate
// Will insert the updated data into the node, will return a boolean for each data entry
func (s *grpcServer) DataUpdate(stream proto.RouteGuide_DataUpdateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// This is data coming in for the first time, set to update other nodes if required
		updated, updatedKey := data.StoreDataInDatabase(&in.Key, &in.Database, &in.Data, false, int(in.SecondaryNodeGroup))

		if err := stream.Send(&proto.Updated{
			UpdatedSuccessfully: updated,
			UpdatedKey:          *updatedKey,
		}); err != nil {
			return err
		}
	}
}

// ReplicaDataUpdate
// Will insert the updated data into the node, will return a boolean for each data entry
func (s *grpcServer) ReplicaDataUpdate(stream proto.RouteGuide_ReplicaDataUpdateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		stored, key := data.StoreDataInDatabase(&in.Key, &in.Database, &in.Data, true, int(in.SecondaryNodeGroup))

		if err := stream.Send(&proto.Updated{
			UpdatedSuccessfully: stored,
			UpdatedKey:          *key,
		}); err != nil {
			return err
		}
	}
}

// (ctx context.Context, point *pb.Point) (*pb.Feature, error) {
// GetData
// Will attempt to get the data from the provided location
func (s *grpcServer) GetData(ctx context.Context, reqData *proto.Data) (*proto.Data, error) {
	exists, discoveredData := data.RetrieveDataFromDatabase(&reqData.Key, &reqData.Database)
	if exists {
		return &proto.Data{
			Key:              reqData.Key,
			Database:         reqData.Database,
			Data:             *discoveredData,
			ReplicaNodeGroup: reqData.ReplicaNodeGroup,
		}, nil
	} else {
		return &proto.Data{
			Key:              reqData.Key,
			Database:         reqData.Database,
			ReplicaNodeGroup: reqData.ReplicaNodeGroup,
			Data:             nil,
		}, nil
	}
}

// DeleteData
// Will attempt to delete the data from the provided location, will return with a boolean for success status
func (s *grpcServer) DeleteData(stream proto.RouteGuide_DeleteDataServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		deleted, err := data.DeleteDataFromDatabase(&in.Key, &in.Database)
		if !deleted || err != nil {
			return err
		}

		if err := stream.Send(&proto.Updated{
			UpdatedSuccessfully: deleted,
			UpdatedKey:          in.Key,
		}); err != nil {
			return err
		}
	}
}
