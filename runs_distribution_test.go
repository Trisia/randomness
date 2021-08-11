package randomness

import (
	"fmt"
	"testing"
)

func TestRunsDistributionTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := RunsDistributionTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
