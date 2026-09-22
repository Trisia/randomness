package detect

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/Trisia/randomness"
)

// chunkReader 按固定上限分块返回给定数据，用于模拟「允许短读」的合法 io.Reader
// （net.Conn、管道、限速源都会这样）。
type chunkReader struct {
	data []byte
	pos  int
	max  int
}

func (r *chunkReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := len(p)
	if r.max > 0 && n > r.max {
		n = r.max
	}
	if n > len(r.data)-r.pos {
		n = len(r.data) - r.pos
	}
	copy(p, r.data[r.pos:r.pos+n])
	r.pos += n
	return n, nil
}

// failAfterReader 在成功返回若干次数据后开始报错，用于验证错误路径不会死锁。
type failAfterReader struct {
	okLeft int
}

func (r *failAfterReader) Read(p []byte) (int, error) {
	if r.okLeft <= 0 {
		return 0, errors.New("随机源已耗尽")
	}
	r.okLeft--
	for i := range p {
		p[i] = 0xA5
	}
	return len(p), nil
}

// Fast 系列在随机源报错时必须返回错误，而不是永久阻塞在 wg.Wait()。
// 旧实现在错误分支 continue，跳过了 wait.Done()，导致死锁。
func TestFastDetectNoHangOnSourceError(t *testing.T) {
	cases := []struct {
		name string
		run  func(io.Reader) (bool, error)
	}{
		{"FactoryDetectFast", FactoryDetectFast},
		{"PowerOnDetectFast", PowerOnDetectFast},
		{"PeriodDetectFast", PeriodDetectFast},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			type res struct {
				pass bool
				err  error
			}
			done := make(chan res, 1)
			go func() {
				pass, err := c.run(&failAfterReader{okLeft: 3})
				done <- res{pass, err}
			}()
			select {
			case got := <-done:
				if got.err == nil {
					t.Fatalf("随机源报错时 %s 应当返回错误，实际返回 (%v, nil)", c.name, got.pass)
				}
				if got.pass {
					t.Fatalf("%s 在报错时不应返回 pass=true", c.name)
				}
				t.Logf("%s 正确返回错误: %v", c.name, got.err)
			case <-time.After(20 * time.Second):
				t.Fatalf("%s 在随机源报错时死锁（wg.Wait 未返回）", c.name)
			}
		})
	}
}

// Fast 版本必须用 io.ReadFull 组装完整样本。
// 用一个每次只给 1 字节的 reader 与整块 reader 对比，结果必须完全一致；
// 旧实现直接把短读的残缺缓冲区送进检测，会得出不同（且错误）的结论。
func TestFastDetectHandlesShortReads(t *testing.T) {
	seed := uint64(20240920)
	// PeriodDetect 需要 20 组 × 2500 字节
	data := randomness.NewDetRand(seed).RawBytes(20 * 2500)

	bulk := &chunkReader{data: data, max: 0}
	short := &chunkReader{data: data, max: 1}
	trickle := &chunkReader{data: data, max: 7}

	passBulk, errBulk := PeriodDetectFast(bulk)
	passShort, errShort := PeriodDetectFast(short)
	passTrickle, errTrickle := PeriodDetectFast(trickle)

	if errBulk != nil || errShort != nil || errTrickle != nil {
		t.Fatalf("短读不应导致错误: bulk=%v short=%v trickle=%v", errBulk, errShort, errTrickle)
	}
	if passBulk != passShort || passBulk != passTrickle {
		t.Fatalf("短读改变了结论: bulk=%v short=%v trickle=%v", passBulk, passShort, passTrickle)
	}
	t.Logf("整块/1字节/7字节 读取结论一致: pass=%v", passBulk)
}

// PeriodDetectFast 必须与串行版 PeriodDetect 语义一致（同为 12 项、同样 20 组 × 2×10^4 bit）。
// 这条同时守住了「Fast 版误用 Round15」的回归。
func TestPeriodDetectFastMatchesSerial(t *testing.T) {
	const total = 20 * 2500
	data := randomness.NewDetRand(20240920).RawBytes(total)

	passSerial, errSerial := PeriodDetect(&chunkReader{data: data})
	passFast, errFast := PeriodDetectFast(&chunkReader{data: data})

	if (errSerial == nil) != (errFast == nil) {
		t.Fatalf("错误不一致: serial=%v fast=%v", errSerial, errFast)
	}
	if errSerial != nil {
		t.Fatalf("两者都报错: serial=%v fast=%v", errSerial, errFast)
	}
	if passSerial != passFast {
		t.Fatalf("PeriodDetect=%v 与 PeriodDetectFast=%v 结论不一致（Fast 版检测项数与串行版不同？）",
			passSerial, passFast)
	}
	t.Logf("PeriodDetect 与 PeriodDetectFast 结论一致: pass=%v", passFast)
}

// 单次检测也应当在短读源上稳定工作。
func TestSingleDetectHandlesShortReads(t *testing.T) {
	data := randomness.NewDetRand(7).RawBytes(4096)
	for _, numByte := range []int{16, 40, 160, 1280} {
		bulk, errBulk := SingleDetect(&chunkReader{data: data}, numByte)
		short, errShort := SingleDetect(&chunkReader{data: data, max: 1}, numByte)
		if (errBulk == nil) != (errShort == nil) {
			t.Fatalf("numByte=%d 错误不一致: bulk=%v short=%v", numByte, errBulk, errShort)
		}
		if errBulk == nil && bulk != short {
			t.Fatalf("numByte=%d 短读改变了结论: bulk=%v short=%v", numByte, bulk, short)
		}
	}
}
