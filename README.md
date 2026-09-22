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

### BitSeq：一次转换，多项复用（推荐）

要对同一段数据跑多项检测时，请先用构造器把字节打包成 `randomness.BitSeq`，
再调用 `XxxTestBitSeq` / `XxxBitSeq` 系列。`BitSeq` 是位序列的紧凑表示
（1 位/比特，即 `[]bool` 的 1/8 内存），也是库内全部 15 种检测的算法内核
实际接受的类型。

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

> **为什么推荐这个用法**：若改用 `B2bitArr` + `XxxProto` / `XxxTest` 的组合，
> 同一段数据会先被展开成 `[]bool`（1 字节/位），随后每项检测内部还要把它
> 重新打包一次——一轮检测里同一份数据被反复转换十余遍。在 10^8 bit 规模下，
> 这相当于常驻 100MB 的展开数组加上十余次的重复转换。

`BitSeq` 提供的公开操作：

| 操作 | 说明 |
|---|---|
| `NewBitSeqFromBytes(data)` / `BitSeqFromBytes(data)` | 从 `[]byte` 构造（构造器风格与函数风格，两者等价） |
| `BitSeqFromBools(bs)` | 从 `[]bool` 构造 |
| `NewBitSeq(nbits)` | 构造长度为 `nbits` 的全 0 序列 |
| `ReadBitSeqFromFile(filename)` | 从文件读取二进制字节构造 |
| `Len()` | 返回序列长度（比特数） |
| `Bit(i)` | 返回第 i 位的值；越界返回 `false` |
| `Set(i, v)` | 设置第 i 位；越界不操作 |
| `Ones()` | 返回全部有效位中 1 的个数 |
| `OnesInRange(lo, hi)` | 返回 `[lo, hi)` 区间内 1 的个数 |
| `ToBools()` | 展开为 `[]bool` |
| `ToBytes()` | 按 MSB 优先打包为 `[]byte` |
| `Copy()` | 深拷贝（独立底层存储） |

位序与 `B2bitArr` 完全一致，因此两者可以自由互换。

需要生成可复现的测试数据时，可使用 `rdgen -seed <非0>`（见其 README）。

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

> 注意：离散傅里叶检测在 10^8 bit 规模下单次检测约需 **2.05 GB** 峰值内存（用 `rddetector` 实测 `Maximum resident set size`）——实数打包后的数据为 `2^26` 个 `complex128`（约 1.0 GB），FFT 根表约 1.0 GB，输入 `BitSeq` 约 12.5 MB。若并发执行多次检测请相应控制并发数，防止内存溢出（OOM）。
>
> 同一份数据若沿用旧的 `B2bitArr` + 全长复数 FFT 路径，峰值约为 **5.1 GB**。
> 10^6 bit 规模下单次约需 8.5 MB。文本此前标注的 "1024MB" 系更早实现的数据，已按实测更正。

## 发展

检测规范：
- [GM/T 0005-2021《随机性检测规范》](http://www.gmbz.org.cn/main/viewfile/20220805030323734119.html)
- GB/T 32915-2016《信息安全技术 二元序列随机性检测方法》
- GM/T 0005-2012《随机性检测规范》 *已废止*

在 GM/T 0005-2021《随机性检测规范》中在 **样本通过率判定** 的基础上，增加了 **样本均匀性判定** 作为检测通过判定依据，详见 GM/T 0005-2021
《随机性检测规范》 6.3。

目前 **randomness** 已经升级至 GM/T 0005-2021《随机性检测规范》。