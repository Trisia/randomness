package detect

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestFactoryDetectFast(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	hit := "通过"
	pass, err := FactoryDetectFast(rand.Reader)
	if err != nil {
		hit = err.Error()
	}
	fmt.Printf("15种算法 上电自检 50组 10^6 bit: %v, hit: %s\n", pass, hit)
}

func TestPowerOnDetectFast(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	hit := "通过"
	pass, err := PowerOnDetectFast(rand.Reader)
	if err != nil {
		hit = err.Error()
	}
	fmt.Printf("15种算法 上电自检 20组 10^6 bit: %v, hit: %s\n", pass, hit)
}

func TestPeriodDetectFast(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	hit := "通过"
	pass, err := PeriodDetectFast(rand.Reader)
	if err != nil {
		hit = err.Error()
	}
	fmt.Printf("12种算法 周期检测 20组 20000 bit: %v, hit: %s\n", pass, hit)
}
