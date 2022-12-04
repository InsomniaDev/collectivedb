package database

import (
	"testing"
)

func TestInsert(t *testing.T) {
	validKey := "testKey"
	validValue := "testValue"

	invalidKey := ""
	invalidValue := ""

	type args struct {
		key   *string
		value *string
	}
	tests := []struct {
		name            string
		args            args
		wantInserted    bool
		wantInsertedKey *string
	}{
		{
			name: "Insert Succeed",
			args: args{
				key:   &validKey,
				value: &validValue,
			},
			wantInserted:    true,
			wantInsertedKey: &validKey,
		},
		{
			name: "Insert Failed",
			args: args{
				key:   &invalidKey,
				value: &invalidValue,
			},
			wantInserted:    false,
			wantInsertedKey: &invalidKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInserted, gotInsertedKey := Insert(tt.args.key, tt.args.value)
			if gotInserted != tt.wantInserted {
				t.Errorf("Insert() gotInserted = %v, want %v", gotInserted, tt.wantInserted)
			}
			if gotInsertedKey != tt.wantInsertedKey {
				t.Errorf("Insert() gotInsertedKey = %v, want %v", gotInsertedKey, tt.wantInsertedKey)
			}
		})
	}
}

func TestEdit(t *testing.T) {
	validKey := "testKey"
	validValue := "testValue"

	invalidKey := ""
	invalidValue := ""

	type args struct {
		key   *string
		value *string
	}
	tests := []struct {
		name           string
		args           args
		wantUpdated    bool
		wantUpdatedKey *string
	}{
		{
			name: "Update Succeed",
			args: args{
				key:   &validKey,
				value: &validValue,
			},
			wantUpdated:    true,
			wantUpdatedKey: &validKey,
		},
		{
			name: "Update Failed",
			args: args{
				key:   &invalidKey,
				value: &invalidValue,
			},
			wantUpdated:    false,
			wantUpdatedKey: &invalidKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUpdated, gotUpdatedKey := Edit(tt.args.key, tt.args.value)
			if gotUpdated != tt.wantUpdated {
				t.Errorf("Edit() gotUpdated = %v, want %v", gotUpdated, tt.wantUpdated)
			}
			if gotUpdatedKey != tt.wantUpdatedKey {
				t.Errorf("Edit() gotUpdatedKey = %v, want %v", gotUpdatedKey, tt.wantUpdatedKey)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	validKey := "testKey"
	invalidKey := ""

	type args struct {
		key *string
	}
	tests := []struct {
		name        string
		args        args
		wantDeleted bool
	}{
		{
			name: "Delete Succeed",
			args: args{
				key:   &validKey,
			},
			wantDeleted:    true,
		},
		{
			name: "Delete Failed",
			args: args{
				key:   &invalidKey,
			},
			wantDeleted:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDeleted := Delete(tt.args.key); gotDeleted != tt.wantDeleted {
				t.Errorf("Delete() = %v, want %v", gotDeleted, tt.wantDeleted)
			}
		})
	}
}
