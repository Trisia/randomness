package randomness

// Runner 运行随机性接口
// bits: 比特序列
// return: P-value 衡量样本随机性好坏的度量指标
type Runner func(bits []bool) float64
