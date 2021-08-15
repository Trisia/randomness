package randomness

import (
	"fmt"
	"testing"
)

func TestRunsTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := RunsTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
