package database

import (
	"reflect"
	"testing"
)

func TestUpdate(t *testing.T) {
	bucket := "test"

	validKey := "testKey"
	validValue := "testValue"

	invalidKey := ""
	invalidValue := ""

	type args struct {
		key    *string
		value  *string
		bucket *string
	}
	tests := []struct {
		name           string
		args           args
		wantUpdated    bool
		wantUpdatedKey *string
	}{
		{
			name: "Success",
			args: args{
				key:    &validKey,
				value:  &validValue,
				bucket: &bucket,
			},
			wantUpdated:    true,
			wantUpdatedKey: &validKey,
		},
		{
			name: "Failure",
			args: args{
				key:    &invalidKey,
				value:  &invalidValue,
				bucket: &bucket,
			},
			wantUpdated:    false,
			wantUpdatedKey: &invalidKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInserted, gotInsertedKey := Update(tt.args.key, tt.args.value, tt.args.bucket)
			if gotInserted != tt.wantUpdated {
				t.Errorf("Insert() gotInserted = %v, want %v", gotInserted, tt.wantUpdated)
			}
			if gotInsertedKey != tt.wantUpdatedKey {
				t.Errorf("Insert() gotInsertedKey = %v, want %v", gotInsertedKey, tt.wantUpdatedKey)
			}
		})
	}
}

func TestGet(t *testing.T) {
	bucket := "test"

	validKey := "testKey"
	validValue := "testValue"
	returnedValue := []byte(validValue)

	invalidKey := ""

	type args struct {
		key    *string
		bucket *string
	}
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantValue  *[]byte
	}{
		{
			name: "Success",
			args: args{
				key:    &validKey,
				bucket: &bucket,
			},
			wantExists: true,
			wantValue:  &returnedValue,
		},
		{
			name: "Failure",
			args: args{
				key:    &invalidKey,
				bucket: &bucket,
			},
			wantExists: false,
			wantValue:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, gotValue := Get(tt.args.key, tt.args.bucket)
			if gotExists != tt.wantExists {
				t.Errorf("Get() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Get() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func Test_getDatabase(t *testing.T) {
	bucket := "test"
	newBucket := "new"

	type args struct {
		bucket *string
	}
	tests := []struct {
		name            string
		args            args
		wantConnections int
	}{
		{
			name: "Retrieved",
			args: args{
				bucket: &bucket,
			},
			wantConnections: 1,
		},
		{
			name: "Create New",
			args: args{
				bucket: &newBucket,
			},
			wantConnections: 2,
		},
		{
			name: "Didn't create new",
			args: args{
				bucket: &newBucket,
			},
			wantConnections: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDatabase(tt.args.bucket)
			if got != nil {
				if len(connections) != tt.wantConnections {
					t.Errorf("Total Connections = %d, wantConnections %d", len(connections), tt.wantConnections)
					return
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	bucket := "test"

	validKey := "testKey"
	invalidKey := ""

	type args struct {
		key    *string
		bucket *string
	}
	tests := []struct {
		name        string
		args        args
		wantDeleted bool
		wantErr     bool
	}{
		{
			name: "Success",
			args: args{
				key:    &validKey,
				bucket: &bucket,
			},
			wantDeleted: true,
			wantErr:     false,
		},
		{
			name: "Failure",
			args: args{
				key:    &invalidKey,
				bucket: &bucket,
			},
			wantDeleted: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeleted, err := Delete(tt.args.key, tt.args.bucket)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDeleted != tt.wantDeleted {
				t.Errorf("Delete() = %v, want %v", gotDeleted, tt.wantDeleted)
			}
		})
	}
}
