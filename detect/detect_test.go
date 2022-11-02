package detect

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestFactoryDetect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	hit := "通过"
	pass, err := FactoryDetect(rand.Reader)
	if err != nil {
		hit = err.Error()
	}
	fmt.Printf("15种算法 上电自检 50组 10^6 bit: %v, hit: %s\n", pass, hit)
}

func TestPowerOnDetect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	hit := "通过"
	pass, err := PowerOnDetect(rand.Reader)
	if err != nil {
		hit = err.Error()
	}
	fmt.Printf("15种算法 上电自检 20组 10^6 bit: %v, hit: %s\n", pass, hit)
}

func TestPeriodDetect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	hit := "通过"
	pass, err := PeriodDetect(rand.Reader)
	if err != nil {
		hit = err.Error()
	}
	fmt.Printf("12种算法 周期检测 20组 20000 bit: %v, hit: %s\n", pass, hit)
}

func TestSingleDetect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	pass, err := SingleDetect(rand.Reader, 16)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("单次检测 320 bit:", pass)
	pass, err = SingleDetect(rand.Reader, 320/8)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("扑克检测 单次检测 10^6 bit:", pass)
}

func TestThresholdQ(t *testing.T) {
	qValues := []float64{0.9, 0.91, 0.07, 0.08, 0.1, 0.11, 0.12, 0.13, 0.14, 0.2, 0.21, 0.22, 0.23, 0.24, 0.25, 0.26, 0.27, 0.3, 0.31, 0.32, 0.33, 0.34, 0.35, 0.36, 0.4, 0.45, 0.5, 0.51, 0.52, 0.53, 0.54, 0.6, 0.61, 0.7, 0.71, 0.72, 0.73, 0.74, 0.75, 0.76, 0.77, 0.8, 0.81, 0.82, 0.83, 0.84, 0.85, 0.86, 0.87, 0.88}
	result := ThresholdQ(qValues)
	if fmt.Sprintf("%.6f", result) != "0.096578" {
		t.FailNow()
	}
}

/*
func TestPowerOnDetect2(t *testing.T) {
	var files []io.Reader
	dirname := "D:\\Project\\cliven\\randomness\\tools\\rdgen\\target\\data"
	dirs, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	defer func() {
		for _, file := range files {
			_ = file.(io.Closer).Close()
		}
	}()
	for _, fi := range dirs {
		name := filepath.Join(dirname, fi.Name())
		file, err := os.OpenFile(name, os.O_RDONLY, os.FileMode(666))
		if err != nil {
			panic(err)
		}
		files = append(files, file)
	}
	mReader := io.MultiReader(files...)
	hit := "通过"

	pass, err := PowerOnDetect(mReader)
	if err != nil {
		hit = err.Error()
	}
	fmt.Printf("15种算法 上电自检 20组 10^6 bit: %v, hit: %s\n", pass, hit)
}
*/
