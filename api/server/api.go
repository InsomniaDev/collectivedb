package server

import (
	context "context"
	"io"
	"sync"

	api "github.com/insomniadev/collective-db/api"
	"github.com/insomniadev/collective-db/internal/control"
)

// Server type for working with the gRPC server
type grpcServer struct {
	api.UnimplementedRouteGuideServer

	dictionary_mu sync.Mutex
}

// Create and return the gRPC server
func NewGrpcServer() *grpcServer {
	s := &grpcServer{}
	return s
}

// ReplicaUpdate receives a stream of data updates and
// returns with a boolean on if updated successfully
func (s *grpcServer) ReplicaUpdate(stream api.RouteGuide_ReplicaUpdateServer) error {
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
		control.ReplicaUpdate(convertDataUpdatesToControlDataUpdate(in))
		s.dictionary_mu.Unlock()

		if err := stream.Send(&api.Updated{
			UpdatedSuccessfully: true,
		}); err != nil {
			return err
		}
	}
}

// SyncDataRequest
// Will send a request to the server to pull in all of the data to the newly joined node
func (s *grpcServer) SyncDataRequest(syncIpAddress *api.SyncIp, stream api.RouteGuide_SyncDataRequestServer) (err error) {
	// Make a channel that the process can yield the discovered data through
	storedData := make(chan *control.StoredData)

	// Setup a process to wait for the returned data
	var wg sync.WaitGroup

	// Go through and return the data as it is discovered
	wg.Add(1)
	go func(chan *control.StoredData) {
		data := <-storedData
		if data != nil {
			if err = stream.Send(&api.Data{
				Key:      data.DataKey,
				Database: data.Database,
				Data:     data.Data,
			}); err != nil {
				// Stop processing since we hit an error
				defer wg.Done()
				return
			}
		} else {
			// We are done processing
			defer wg.Done()
			return
		}
	}(storedData)

	// Request the data to be returned
	control.ReplicaSyncRequest(storedData)

	// Wait until all stream data has been sent
	wg.Wait()
	return err
}

// DictionaryUpdate
// Will send a stream of data entries that requie an update, will respond with a boolean for each entry sent
func (s *grpcServer) DictionaryUpdate(stream api.RouteGuide_DictionaryUpdateServer) error {
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
		control.CollectiveUpdate(convertDataUpdatesToControlDataUpdate(in))
		s.dictionary_mu.Unlock()

		if err := stream.Send(&api.Updated{
			UpdatedSuccessfully: true,
		}); err != nil {
			return err
		}
	}
}

// DataUpdate
// Will insert the updated data into the node, will return a boolean for each data entry
func (s *grpcServer) DataUpdate(stream api.RouteGuide_DataUpdateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		updated, updatedKey := control.ReplicaStoreData(in.Key, in.Database, in.Data)

		if err := stream.Send(&api.Updated{
			UpdatedSuccessfully: updated,
			UpdatedKey:          *updatedKey,
		}); err != nil {
			return err
		}
	}
}

// ReplicaDataUpdate
// Will insert the updated data into the node, will return a boolean for each data entry
func (s *grpcServer) ReplicaDataUpdate(stream api.RouteGuide_ReplicaDataUpdateServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		stored, key := control.ReplicaStoreData(in.Key, in.Database, in.Data)

		if err := stream.Send(&api.Updated{
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
func (s *grpcServer) GetData(ctx context.Context, data *api.Data) (*api.Data, error) {
	exists, discoveredData := control.RetrieveData(&data.Key, &data.Database)
	if exists {
		return &api.Data{
			Key:      data.Key,
			Database: data.Database,
			Data:     *discoveredData,
		}, nil
	} else {
		return &api.Data{
			Key:      data.Key,
			Database: data.Database,
			Data:     nil,
		}, nil
	}
}

// DeleteData
// Will attempt to delete the data from the provided location, will return with a boolean for success status
func (s *grpcServer) DeleteData(ctx context.Context, data *api.Data) (*api.Updated, error) {
	deleted, err := control.DeleteData(&data.Key, &data.Database)
	return &api.Updated{
		UpdatedSuccessfully: deleted,
		UpdatedKey:          data.Key,
	}, err
}
