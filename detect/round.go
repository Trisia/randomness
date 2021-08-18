package detect

import "github.com/Trisia/randomness"

// Round15 15种方法测试轮
// data: 待检测数据，推荐长度： 10^6 bit =>  125,000 byte
func Round15(data []byte) []*randomness.TestResult {
	results := make([]*randomness.TestResult, 15)
	for i, method := range randomness.TestMethodArr {
		results[i] = method.Runner(data)
	}
	return results
}

// Round12 12种方法测试轮（除去：离散傅里叶检测、线型复杂度检测、通用统计）
// data: 待检测数据，推荐长度： 20000 bit =>  2,500
func Round12(data []byte) []*randomness.TestResult {
	results := make([]*randomness.TestResult, 12)
	arr := randomness.TestMethodArr[:12]
	for i, method := range arr {
		results[i] = method.Runner(data)
	}
	return results
}
