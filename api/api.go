package api

import (
	context "context"
	"io"
	"sync"

	"github.com/insomniadev/collective-db/internal/control"
)

// Server type for working with the gRPC server
type grpcServer struct {
	UnimplementedRouteGuideServer

	dictionary_mu sync.Mutex
}

// Create and return the gRPC server
func NewGrpcServer() *grpcServer {
	s := &grpcServer{}
	return s
}

// ReplicaUpdate receives a stream of data updates and
// returns with a boolean on if updated successfully
func (s *grpcServer) ReplicaUpdate(stream RouteGuide_ReplicaUpdateServer) error {
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

		if err := stream.Send(&Updated{
			UpdatedSuccessfully: true,
		}); err != nil {
			return err
		}
	}
}

// SyncDataRequest
// Will send a request to the server to pull in all of the data to the newly joined node
func (s *grpcServer) SyncDataRequest(syncIpAddress *SyncIp, stream RouteGuide_SyncDataRequestServer) error {
	return nil
}

// DictionaryUpdate
// Will send a stream of data entries that requie an update, will respond with a boolean for each entry sent
func (s *grpcServer) DictionaryUpdate(stream RouteGuide_DictionaryUpdateServer) error {
	return nil
}

// DataUpdate
// Will insert the updated data into the node, will return a boolean for each data entry
func (s *grpcServer) DataUpdate(stream RouteGuide_DataUpdateServer) error {
	return nil
}

// ReplicaDataUpdate
// Will insert the updated data into the node, will return a boolean for each data entry
func (s *grpcServer) ReplicaDataUpdate(stream RouteGuide_ReplicaDataUpdateServer) error {
	return nil
}

// (ctx context.Context, point *pb.Point) (*pb.Feature, error) {
// GetData
// Will attempt to get the data from the provided location
func (s *grpcServer) GetData(ctx context.Context, data *Data) (*Data, error) {
	return nil, nil
}

// DeleteData
// Will attempt to delete the data from the provided location, will return with a boolean for success status
func (s *grpcServer) DeleteData(ctx context.Context, data *Data) (*Updated, error) {
	return nil, nil
}
