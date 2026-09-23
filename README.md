# Randomness 二元序列随机性检测方法

[![Documentation](https://godoc.org/github.com/Trisia/randomness?status.svg)](https://pkg.go.dev/github.com/Trisia/randomness) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Trisia/randomness) ![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/Trisia/randomness)

**在使用randomness前，请务必悉知 [*《randomness 免责声明》*](免责声明.md)！**

> - 致谢 [zhebinhu/randomnessTests](https://github.com/zhebinhu/randomnessTests) 项目。
> - 致谢 [Sun Yimin](https://github.com/emmansun) 对本项目关键性建设。

## 项目概述

该工具库实现了《GM/T 0005-2021 随机性检测规范》中描述的15种随机性检测方法。

![rundesp.png](tools/rddetector/rundesp.png)

- [ 1] 单比特频数检测      [MonoBitFrequencyTest](./mono_bit_frequency.go)
- [ 2] 块内频数检测        [FrequencyWithinBlockTest](./frequency_within_block.go)
- [ 3] 扑克检测            [PokerTest](./poker.go)
- [ 4] 重叠子序列检测      [OverlappingTemplateMatchingTest](./overlapping.go)
- [ 5] 游程总数检测        [RunsTest](./runs.go)
- [ 6] 游程分布检测        [RunsDistributionTest](./runs_distribution.go)
- [ 7] 块内最大游程检测    [LongestRunOfOnesInABlockTest](./longest_run_of_ones_In_block.go)
- [ 8] 二元推导检测        [BinaryDerivativeTest](./binary_derivative.go)
- [ 9] 自相关检测          [AutocorrelationTest](./autocorrelation.go)
- [10] 矩阵秩检测          [MatrixRankTest](./matrix_rank.go)
- [11] 累加和检测          [CumulativeTest](./cumulative.go)
- [12] 近似熵检测          [ApproximateEntropyTest](./approximate_entropy.go)
- [13] 线型复杂度检测      [LinearComplexityTest](./linear_complexity.go)
- [14] Maurer通用统计检测  [MaurerUniversalTest](./maurers_universal.go)
- [15] 离散傅里叶检测      [DiscreteFourierTransformTest](./discrete_fourier_transform.go)

### 检测工具

若您需要使用相关测试工具，可以 **[从 Release 中下载编译版本](https://github.com/Trisia/randomness/releases)** 或 **手动编译程序**，详见文档：

- [随机性检测工具 使用说明 rddetector](./tools/rddetector/README.md)
- [数据生成工具 使用说明 rdgen](./tools/rdgen/README.md)

## 快速入门

`rddetector` 检测工具只是 `randomness` 的一种应用方式，`randomness` 提供丰富的 API 接口，您可以根据需要定制化使用。

安装 `randomness`：

```bash
go get -u github.com/Trisia/randomness
```

下面是 $10^6$ 比特数据规模的扑克检测的例子：

```go
package main

import (
	"crypto/rand"
	"fmt"
	"github.com/Trisia/randomness"
)

func main() {
	// 产生随机数序列
	n := 1000000
	buf := make([]byte, n/8)
	_, _ = rand.Read(buf)

	// 直接对字节做检测，内部会自动完成转换
	p, _ := randomness.PokerTestBytes(buf, 8)
	fmt.Printf("扑克检测 n: 1000000, P-value: %f\n", p)
}
```

### 数据检测建议

推荐使用 `BitSeq` 后缀的接口来实现同一组数据的多次检测，可以有效的减少数据转换和内存占用。（例如：`RunsTestBytes` 对应 `RunsTestBitSeq`）

```go
package main

import (
	"crypto/rand"
	"fmt"
	"github.com/Trisia/randomness"
)

func main() {
	n := 1000000
	buf := make([]byte, n/8)
	_, _ = rand.Read(buf)

	// 使用构造器把字节打包成 BitSeq（10^6 bit 约 0.07ms / 131KB）
	s := randomness.NewBitSeqFromBytes(buf)

	// 后续所有检测复用同一个 BitSeq，不再产生任何转换开销
	p, q := randomness.FrequencyWithinBlockTestBitSeq(s, 10000)
	fmt.Printf("块内频数检测 P: %f Q: %f\n", p, q)

	p, q = randomness.RunsTestBitSeq(s)
	fmt.Printf("游程总数检测 P: %f Q: %f\n", p, q)

	p1, p2, _, _ := randomness.OverlappingTemplateMatchingTestBitSeq(s, 5)
	fmt.Printf("重叠子序列检测 P1: %f P2: %f\n", p1, p2)

	// 或者直接遍历与 TestMethodArr 一一对应的 TestMethodSeqArr
	for _, item := range randomness.TestMethodSeqArr {
		res := item.Runner(s)
		fmt.Printf("%s P: %f Q: %f 通过: %v\n", res.Name, res.P, res.Q, res.Pass)
	}
}
```


`BitSeq`提供了多种方法，您可以实现多种数据类型自由转换和数据的随机访问。

更多 API 使用方法见：[randomness API 文档](https://pkg.go.dev/github.com/Trisia/randomness)

## 随机数发生器检测

`randomness`
实现了 [GM/T 0062-2018《密码产品随机数检测要求》](http://www.gmbz.org.cn/main/viewfile/20180519123243999421.html) 9. E类产品随机数检测
随机数发生器 4个不同应用阶段的随机数检测：

- 出厂检测 [detect.FactoryDetect](detect/detect.go)
- 上电检测 [detect.PowerOnDetect](detect/detect.go)
- 周期检测 [detect.PeriodDetect](detect/detect.go)
- 单次检测 [detect.SingleDetect](detect/detect.go)

使用方法见 [测试用例 detect_test.go](detect/detect_test.go)

如果您的主机处理器含有多个核心，那么可以使用 Fast 系列的 API 来加速检测，见 [测试用例 detect_fast_test.go](detect/detect_fast_test.go)

## 性能参数

> 环境：Go 1.25.7 / 16 核 CPU；数据为 `crypto/rand` 生成的随机比特序列。
> 测量方式：每项检测预热 1 次后重复测量取均值（2×10⁴ bit 重复 20 次、10⁶ bit 重复 10 次、10⁸ bit 重复 1 次）。

### 耗时

单位：ms

| 检测项 | 2×10⁴ bit | 10⁶ bit | 10⁸ bit |
|:-------|--------:|--------:|--------:|
| 单比特频数检测 | 0.003 | 0.093 | 7.918 |
| 块内频数检测 | 0.004 | 0.093 | 8.132 |
| 扑克检测 | 0.004 | 0.160 | 7.082 |
| 重叠子序列检测 | 0.038 | 1.665 | 129.981 |
| 游程总数检测 | 0.004 | 0.138 | 12.710 |
| 游程分布检测 | 0.023 | 0.929 | 91.745 |
| 块内最大“1”游程检测 | 0.014 | 0.442 | 54.482 |
| 二元推导检测 | 0.006 | 0.161 | 21.655 |
| 自相关检测 | 0.005 | 0.099 | 10.166 |
| 矩阵秩检测 | 0.121 | 5.066 | 492.435 |
| 累加和检测 | 0.029 | 1.306 | 127.904 |
| 近似熵检测 | 0.068 | 2.837 | 274.047 |
| 线型复杂度检测 | 0.700 | 7.672 | 494.494 |
| 通用统计检测 | 0.071 | 5.378 | 485.173 |
| 离散傅里叶检测 | 0.547 | 42.011 | 12827.069 |

### 内存

单位：KB（单次调用累计分配量，等价 `go test -bench` 的 B/op）

| 检测项 | 2×10⁴ bit | 10⁶ bit | 10⁸ bit |
|:-------|--------:|--------:|--------:|
| 单比特频数检测 | 0.06 | 0.06 | 0.06 |
| 块内频数检测 | 2.69 | 128.06 | 12208.06 |
| 扑克检测 | 2.06 | 2.06 | 2.06 |
| 重叠子序列检测 | 3.12 | 128.50 | 12208.50 |
| 游程总数检测 | 2.69 | 128.06 | 12208.06 |
| 游程分布检测 | 2.92 | 128.44 | 12208.58 |
| 块内最大“1”游程检测 | 2.73 | 129.38 | 12209.38 |
| 二元推导检测 | 5.31 | 256.06 | 24416.16 |
| 自相关检测 | 2.69 | 128.06 | 12208.06 |
| 矩阵秩检测 | 2.94 | 128.31 | 12208.31 |
| 累加和检测 | 2.69 | 128.06 | 12208.06 |
| 近似熵检测 | 3.44 | 128.81 | 12208.81 |
| 线型复杂度检测 | 3.02 | 142.30 | 12223.03 |
| 通用统计检测 | 3.69 | 129.06 | 12209.06 |
| 离散傅里叶检测 | 258.69 | 8320.06 | 1060784.06 |

> 10⁶ bit 下离散傅里叶检测最慢（约 42 ms），是整轮耗时的主要来源，其余单项均在 8 ms 以内；
> 内存方面，除离散傅里叶检测外，其余单项单次分配普遍在 130 KB 左右。
> 完整一轮（15 项）在三种规模下的耗时约为 **1.6 ms / 61 ms / 15.8 s**。

## 发展

检测规范：
- [GM/T 0005-2021《随机性检测规范》](http://www.gmbz.org.cn/main/viewfile/20220805030323734119.html)
- GB/T 32915-2016《信息安全技术 二元序列随机性检测方法》
- GM/T 0005-2012《随机性检测规范》 *已废止*

在 GM/T 0005-2021《随机性检测规范》中在 **样本通过率判定** 的基础上，增加了 **样本均匀性判定** 作为检测通过判定依据，详见 GM/T 0005-2021
《随机性检测规范》 6.3。

目前 **randomness** 已经升级至 GM/T 0005-2021《随机性检测规范》。
