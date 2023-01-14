package control

import (
	"reflect"
	"testing"
)

func Test_storeDataInDatabase(t *testing.T) {
	key := "test"
	bucket := "test"
	data := []byte("hello")

	emptyKey := ""

	type args struct {
		key          *string
		bucket       *string
		data         *[]byte
		replicaStore bool
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *string
	}{
		{
			name: "Store without key",
			args: args{
				key:          &emptyKey,
				bucket:       &bucket,
				data:         &data,
				replicaStore: false,
			},
			want:  true,
			want1: &key,
		},
		{
			name: "Stored successfully",
			args: args{
				key:          &key,
				bucket:       &bucket,
				data:         &data,
				replicaStore: false,
			},
			want:  true,
			want1: &key,
		},
		{
			name: "Stored replica successfully",
			args: args{
				key:          &key,
				bucket:       &bucket,
				data:         &data,
				replicaStore: true,
			},
			want:  true,
			want1: &key,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := storeDataInDatabase(tt.args.key, tt.args.bucket, tt.args.data, tt.args.replicaStore)
			if got != tt.want {
				t.Errorf("storeDataInDatabase() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				if tt.name == "Store without key" {
					if len(*got1) != len("1110b39e-14fd-4a20-b6ed-199709b14eac") {
						t.Errorf("storeDataInDatabase() got1 = %v, want %v", got1, tt.want1)
					}
				} else {
					t.Errorf("storeDataInDatabase() got1 = %v, want %v", got1, tt.want1)
				}
			}
		})
	}
}

func Test_retrieveDataFromDatabase(t *testing.T) {
	key := "test"
	bucket := "test"
	data := []byte("hello")
	storeDataInDatabase(&key, &bucket, &data, false)

	wrongKey := "don'texist"

	type args struct {
		key    *string
		bucket *string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *[]byte
	}{
		{
			name: "Retrieve Successfully",
			args: args{
				key:    &key,
				bucket: &bucket,
			},
			want:  true,
			want1: &data,
		},
		{
			name: "Doesnt exist",
			args: args{
				key:    &wrongKey,
				bucket: &bucket,
			},
			want:  false,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := retrieveDataFromDatabase(tt.args.key, tt.args.bucket)
			if got != tt.want {
				t.Errorf("retrieveDataFromDatabase() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("retrieveDataFromDatabase() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_deleteDataFromDatabase(t *testing.T) {
	key := "test"
	bucket := "test"
	data := []byte("hello")
	storeDataInDatabase(&key, &bucket, &data, false)

	wrongKey := "don'texist"

	type args struct {
		key    *string
		bucket *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Deleted",
			args: args{
				key: &key,
				bucket: &bucket,
			},
			want: true,
			wantErr: false,
		},
		{
			name: "Doesn't exist",
			args: args{
				key: &wrongKey,
				bucket: &bucket,
			},
			want: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deleteDataFromDatabase(tt.args.key, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteDataFromDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("deleteDataFromDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}
