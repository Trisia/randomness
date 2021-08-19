package detect

import (
	"errors"
	"fmt"
	"github.com/Trisia/randomness"
	"io"
	"math"
)

// FactoryDetect 出厂检测，15种检测，每组 10^6比特，分50组
// source: 随机源
func FactoryDetect(source io.Reader) (bool, error) {
	s := 50
	t := Threshold(s)
	buf := make([]byte, 1000_000/8)
	counters := make([]int, 15)
	for i := 0; i < s; i++ {
		_, err := source.Read(buf)
		if err != nil {
			return false, err
		}
		resArr := Round15(buf)
		for idx, result := range resArr {
			if result.Pass {
				counters[idx]++
			}
		}
	}
	for i, n := range counters {
		if n < t {
			return false, fmt.Errorf("%s %d/%d", randomness.TestMethodArr[i].Name, n, s)
		}
	}
	return true, nil
}

// PowerOnDetect 上电自检，15种检测，每组 10^6比特，分20组
// source: 随机源
func PowerOnDetect(source io.Reader) (bool, error) {
	s := 20
	t := Threshold(s)
	buf := make([]byte, 1000_000/8)
	counters := make([]int, 15)
	for i := 0; i < s; i++ {
		_, err := source.Read(buf)
		if err != nil {
			return false, err
		}
		resArr := Round15(buf)
		for idx, result := range resArr {
			if result.Pass {
				counters[idx]++
			}
		}
	}
	for i, n := range counters {
		if n < t {
			return false, fmt.Errorf("%s %d/%d", randomness.TestMethodArr[i].Name, n, s)
		}
	}
	return true, nil
}

// PeriodDetect 周期性检测，除去离散傅里叶检测、线型复杂度检测、通用统计的12种检测
// 检测 20组，每组 20000比特
// source: 随机源
func PeriodDetect(source io.Reader) (bool, error) {
	s := 20
	t := Threshold(s)
	buf := make([]byte, 20000/8)
	counters := make([]int, 12)
	for i := 0; i < s; i++ {
		_, err := source.Read(buf)
		if err != nil {
			return false, err
		}
		resArr := Round12(buf)
		for idx, result := range resArr {
			if result.Pass {
				counters[idx]++
			}
		}
	}
	for i, n := range counters {
		if n < t {
			return false, fmt.Errorf("%s %d/%d", randomness.TestMethodArr[i].Name, n, s)
		}
	}
	return true, nil
}

// SingleDetect 单次检测，单根据实际应用时每次才随机数的大小确定，检测采用扑克检测
// source: 随机源
// numByte: 采集字节数，不能小于16
func SingleDetect(source io.Reader, numByte int) (bool, error) {
	data := make([]byte, numByte)
	_, err := source.Read(data)
	if err != nil {
		return false, err
	}
	n := len(data) * 8
	if n < 128 {
		return false, errors.New("长度不能低于128比特（16 byte）")
	}
	m := 4
	if n < 320 {
		m = 2
	} else if n/8 >= 1280 { // n/m >= 5 * 2^m
		m = 8
	}
	p := randomness.PokerTestBytes(data, m)
	return p >= randomness.Alpha, nil
}

// Threshold 样本通过检测判定数量
// s: 检测样本数
// return 通过检测需要的样本数量
func Threshold(s int) int {
	a := randomness.Alpha
	_s := float64(s)
	r := _s * (1 - a - 3*math.Sqrt((a*(1-a))/_s))

	return int(math.Ceil(r))
}
