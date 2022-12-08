package control

import (
	"testing"
)

func TestStoreData(t *testing.T) {
	bucket := "test"

	testKey := "key"
	testValue := []byte("value")

	type args struct {
		key    *string
		bucket *string
		data   *[]byte
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *string
	}{
		// TODO: Add test cases.
		{
			name: "Stored",
			args: args{
				key:    &testKey,
				bucket: &bucket,
				data:   &testValue,
			},
			want: false,
			want1: &testKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := StoreData(tt.args.key, tt.args.bucket, tt.args.data)
			if got != tt.want {
				t.Errorf("StoreData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("StoreData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
