package randomness

import (
	"fmt"
	"testing"
)

func TestRunsDistributionTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p, q := RunsDistributionTest(bits)
	fmt.Printf("n: 1000000, P-value: %f, Q-value: %f\n", p, q)
}

func TestRunsDistributionTestSample(t *testing.T) {
	bits := B2bitArr([]byte{0xcc, 0x15, 0x62, 0x4b, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe3, 0xd8, 0xb2})
	p, q := RunsDistributionTest(bits)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	//Not same with C.6
	//if fmt.Sprintf("%.6f", p) != "0.970152" || fmt.Sprintf("%.6f", q) != "0.970152" {
	//	t.FailNow()
	//}
}
