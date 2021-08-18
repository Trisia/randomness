# Randomness 二元序列随机性检测方法

![-](https://img.shields.io/badge/language-go-blue.svg) [![license](https://img.shields.io/badge/license-Mulan--2.0-blue)](./LICENSE)


> 感谢 zhebinhu 项目对基础算法实现
> 
> - [zhebinhu/randomnessTests](https://github.com/zhebinhu/randomnessTests)

## 快速入门

运行下面命令安装检测工具集合：

```bash
go get -u github.com/Trisia/randomness
```

```go
func main() {
	// 产生随机数序列
	n := 1000_000
	buf := make([]byte, n/8)
	_, _ = rand.Read(buf)
	// 转换为字节数组
	bits := randomness.B2bitArr(buf)

	// 运行测试组
	p := randomness.PokerTest(bits)
	fmt.Printf("扑克检测 n: 1000000, P-value: %f\n", p)
}
```

产品检测方法：

- 出厂检测 [detect.FactoryDetect](detect/detect.go)
- 上电检测 [detect.PowerOnDetect](detect/detect.go)
- 周期检测 [detect.PeriodDetect](detect/detect.go)
- 单次检测 [detect.SingleDetect](detect/detect.go)

使用方法见 [测试用例](detect/detect_test.go)

## 概述

目前《GM/T 0005-2012 随机性检测规范》已经升级至 《GB/T 32915-2016 信息安全技术 二元序列随机性检测方法》。

该工具库实现了《GB/T 32915-2016 信息安全技术 二元序列随机性检测方法》中描述的15中随机性检测方法：

- [ 1] 单比特频数检测      [MonoBitFrequencyTest](./mono_bit_frequency.go)
- [ 2] 块内频数检测        [FrequencyWithinBlockTest](./frequency_within_block.go)
- [ 3] 扑克检测           [PokerTest](./poker.go)
- [ 4] 重叠子序列检测      [OverlappingTemplateMatchingTest](./overlapping.go)
- [ 5] 游程总数检测        [RunsTest](./runs.go)
- [ 6] 游程分布检测        [RunsDistributionTest](./runs_distribution.go)
- [ 7] 块内最大“1”游程检测 [LongestRunOfOnesInABlockTest](./longest_run_of_ones_In_block.go)
- [ 8] 二元推导检测       [BinaryDerivativeTest](./binary_derivative.go)
- [ 9] 自相关检测         [AutocorrelationTest](./autocorrelation.go)
- [10] 矩阵秩检测        [MatrixRankTest](./matrix_rank.go)
- [11] 累加和检测        [CumulativeTest](./cumulative.go)
- [12] 近似熵检测        [ApproximateEntropyTest](./approximate_entropy.go)
- [13] 线型复杂度检测     [LinearComplexityTest](./linear_complexity.go)
- [14] 通用统计检测       [MaurerUniversalTest](./maurers_universal.go)
- [15] 离散傅里叶检测     [DiscreteFourierTransformTest](./discrete_fourier_transform.go)

检测参数说明：

- 样本长度：`10^6`比特
- 显著水平：`α = 0.01`

各算法检测参数如下：

 序号 | 检测项目 | 参数 |
 --- | --- | --- |
 1 | 单比特频数检测     | - |
 2 | 块内频数检测      | m = 100 |
 3 | 扑克检测        | m = 8 |
 4 | 重叠子序列检测     | m = 5 |
 5 | 游程总数检测      | - |
 6 | 游程分布检测      | - |
 7 | 块内最大“1”游程检测 | m = 10000 |
 8 | 二元推导检测      | k = 7 |
 9 | 自相关检测       | d = 16 |
 10 |  矩阵秩检测      | M = 32, Q= 32 |
 11 |  累加和检测      | - |
 12 |  近似熵检测      | m = 5 |
 13 |  线型复杂度检测    | m = 500 |
 14 |  通用统计检测     | - |
 15 |  离散傅里叶检测    | - |

## 多平台编译

由于Golang支持多平台的跨平台编译，因此可以改造该工具库可以通过Go编译适应平台运行的动态库。

具体方法参照 [《Go语言高级编程》 第2章 CGO编程](https://chai2010.cn/advanced-go-programming-book/)

