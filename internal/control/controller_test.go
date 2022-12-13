package control

import (
	"reflect"
	"testing"
)

func TestStoreData(t *testing.T) {
	bucket := "test"

	newData := ""
	testKey := "key"
	testValue := []byte("value")

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
		// TODO: Add test cases.
		{
			name: "New",
			args: args{
				key:          &newData,
				bucket:       &bucket,
				data:         &testValue,
				replicaStore: false,
			},
			want:  true,
			want1: &newData,
		},
		{
			name: "UpdatedOtherNode",
			args: args{
				key:          &testKey,
				bucket:       &bucket,
				data:         &testValue,
				replicaStore: false,
			},
			want:  true,
			want1: &testKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StoreData(tt.args.key, tt.args.bucket, tt.args.data, tt.args.replicaStore)
			if got != tt.want {
				t.Errorf("StoreData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 && tt.name != "New" {
				t.Errorf("StoreData() got1 = %v, want %v", got1, tt.want1)
			} else {
				if len(*got1) != len("6e79bdbb-1c82-49de-a0ad-abafde999ebc") {
					t.Errorf("StoreData() got1 = %v, want %v", got1, tt.want1)
				}
			}
		})
	}
}

func TestRetrieveData(t *testing.T) {
	bucket := "test"

	testKey := "key"
	testValue := []byte("value")

	testFailKey := "nope"

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
			name: "Found",
			args: args{
				key:    &testKey,
				bucket: &bucket,
			},
			want:  true,
			want1: &testValue,
		},
		{
			name: "Doesn't Exist",
			args: args{
				key:    &testFailKey,
				bucket: &bucket,
			},
			want:  false,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RetrieveData(tt.args.key, tt.args.bucket)
			if got != tt.want {
				t.Errorf("RetrieveData() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("RetrieveData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	bucket := "test"
	testFailBucket := "nope"

	testKey := "key"

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
			name: "Success",
			args: args{
				key:    &testKey,
				bucket: &bucket,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Failure",
			args: args{
				key:    &testKey,
				bucket: &testFailBucket,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteData(tt.args.key, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeleteData() = %v, want %v", got, tt.want)
			}
		})
	}
}
