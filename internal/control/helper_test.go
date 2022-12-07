package control

import (
	"testing"
)

func Test_createUuid(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Created",
			want: "9136b94f-552e-42d8-a7bc-0c5b8acf50df",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createUuid(); len(got) != len(tt.want) {
				t.Errorf("createUuid() = len(%d), want len(%d)", len(got), len(tt.want))
			}
		})
	}
}

func Test_determineIpAddress(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineIpAddress(); got != tt.want {
				t.Errorf("determineIpAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
