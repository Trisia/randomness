package randomness

// Alpha 显著性水平α
const Alpha = 0.01

// AlphaT 分布均匀性的显著性水平
const AlphaT float64 = 0.0001

// TestResult 检测结果
type TestResult struct {
	Name string  // 检测名称
	P    float64 // 检测结果P_value1
	Q    float64 // 检测结果Q_value1
	P2   float64 // 检测结果P_value2
	Q2   float64 // 检测结果Q_value2
	Pass bool    // 是否大于等于显著水平
}

// TestFunc 测试方法
type TestFunc func([]byte) *TestResult

// TestItem 测试项目
type TestItem struct {
	Name string // 检测名称
	// 检测方法
	Runner TestFunc
}

// TestMethodArr 测试方法序列
var TestMethodArr = []TestItem{
	{"单比特频数检测", MonoBitFrequency},
	{"块内频数检测", FrequencyWithinBlock},
	{"扑克检测", Poker},
	{"重叠子序列检测", OverlappingTemplateMatching},
	{"游程总数检测", Runs},
	{"游程分布检测", RunsDistribution},
	{"块内最大“1”游程检测", LongestRunOfOnesInABlock},
	{"二元推导检测", BinaryDerivative},
	{"自相关检测", Autocorrelation},
	{"矩阵秩检测", MatrixRank},
	{"累加和检测", Cumulative},
	{"近似熵检测", ApproximateEntropy},
	{"线型复杂度检测", LinearComplexity},
	{"通用统计检测", MaurerUniversal},
	{"离散傅里叶检测", DiscreteFourierTransform},
}
