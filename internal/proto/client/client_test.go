package client

import (
	"testing"

	"github.com/insomniadev/collectivedb/internal/types"
)

func TestSyncCollectiveRequest(t *testing.T) {
	localhostAddress := "localhost:9091"

	type args struct {
		ipAddress *string
		data      chan<- *types.DataUpdate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "First Test",
			args: args{
				ipAddress: &localhostAddress,
				data: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncCollectiveRequest(tt.args.ipAddress, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("SyncCollectiveRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
