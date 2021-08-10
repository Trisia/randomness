package randomness

import (
	"reflect"
	"testing"
)

func TestB2bit(t *testing.T) {

	tests := []struct {
		name string
		arg  byte
		want []bool
	}{
		{"单字节 1", 0x01, []bool{false, false, false, false, false, false, false, true}},
		{"单字节 2", 0x02, []bool{false, false, false, false, false, false, true, false}},
		{"单字节 4", 0x04, []bool{false, false, false, false, false, true, false, false}},
		{"单字节 8", 0x08, []bool{false, false, false, false, true, false, false, false}},
		{"单字节 16", 0x0C, []bool{false, false, false, false, true, true, false, false}},
		{"单字节 255", 0xFF, []bool{true, true, true, true, true, true, true, true}},
		{"单字节 45", 0x2D, []bool{false, false, true, false, true, true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := B2bit(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("B2bit() = %v, want %v", got, tt.want)
			}
		})
	}
}
