package control

import "testing"

func Test_storeDataInDatabase(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := storeDataInDatabase(tt.args.key, tt.args.bucket, tt.args.data, tt.args.replicaStore)
			if got != tt.want {
				t.Errorf("storeDataInDatabase() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("storeDataInDatabase() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
