package api

import (
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
