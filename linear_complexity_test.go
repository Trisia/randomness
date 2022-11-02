package randomness

import (
	"fmt"
	"sync"
	"testing"
)

var initonce sync.Once
var eConstantBits []bool

func getEConstantBits() []bool {
	initonce.Do(func() {
		eConstantBits = ReadGroupInASCIIFormat("data/data_e")
	})
	return eConstantBits
}

func TestLinearComplexityTestSample(t *testing.T) {
	bits := getEConstantBits()
	p, q := LinearComplexityProto(bits, 1000)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.844721" || fmt.Sprintf("%.6f", q) != "0.844721" {
		t.FailNow()
	}
}
