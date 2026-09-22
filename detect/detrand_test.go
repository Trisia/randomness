package detect

import (
	"testing"

	"github.com/Trisia/randomness"
)

// 用确定性随机源驱动整条检测流程，验证「同一颗种子 -> 同一结论」。
//
// 这一层测的不是某个算法的数值，而是 detect 包的判定逻辑
// （样本通过率 Threshold + 样本均匀性 ThresholdQ）在固定输入下的稳定性，
// 同时也证明 randomness.DetRand 能直接当作 io.Reader 使用。
const detrandSeed = 20240920

func TestSingleDetectDeterministic(t *testing.T) {
	for _, numByte := range []int{16, 40, 128, 160, 1280} {
		ok1, err1 := SingleDetect(randomness.NewDetRand(detrandSeed), numByte)
		ok2, err2 := SingleDetect(randomness.NewDetRand(detrandSeed), numByte)
		if ok1 != ok2 || (err1 == nil) != (err2 == nil) {
			t.Fatalf("numByte=%d 同一种子两次结果不一致: (%v,%v) vs (%v,%v)", numByte, ok1, err1, ok2, err2)
		}
		if err1 != nil {
			t.Logf("SingleDetect(numByte=%d) err=%v", numByte, err1)
		} else {
			t.Logf("SingleDetect(numByte=%d) pass=%v", numByte, ok1)
		}
	}
}

// PeriodDetect 的规模是 20 组 × 2×10⁴ bit，代价很小，适合作为集成级的确定性回归。
func TestPeriodDetectDeterministic(t *testing.T) {
	ok1, err1 := PeriodDetect(randomness.NewDetRand(detrandSeed))
	ok2, err2 := PeriodDetect(randomness.NewDetRand(detrandSeed))
	if ok1 != ok2 || (err1 == nil) != (err2 == nil) {
		t.Fatalf("同一种子两次结果不一致: (%v,%v) vs (%v,%v)", ok1, err1, ok2, err2)
	}
	t.Logf("PeriodDetect(seed=%d) pass=%v err=%v", detrandSeed, ok1, err1)
	// 冻结结论：该种子下 12 项检测应当全部通过。
	// 任一算法的 P 值偏移到足以翻转判定，这里就会失败。
	if !ok1 {
		t.Errorf("PeriodDetect(seed=%d) 由通过变为不通过，err=%v", detrandSeed, err1)
	}
	if err1 != nil {
		t.Errorf("PeriodDetect(seed=%d) 意外报错: %v", detrandSeed, err1)
	}
}

// 保证 DetRand 生成的字节流与 randomness 的位序约定一致。
func TestDetrandFeedsRandomness(t *testing.T) {
	src := randomness.NewDetRand(1)
	data := src.Bytes(20000)
	if len(data) != 20000/8 {
		t.Fatalf("长度 %d，期望 %d", len(data), 20000/8)
	}
	p, q := randomness.MonoBitFrequencyTestBytes(data)
	if p < 0 || p > 1 || q < 0 || q > 1 {
		t.Fatalf("P/Q 越界: %v %v", p, q)
	}
}
