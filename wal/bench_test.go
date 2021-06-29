package wal

import (
	"bytes"
	"testing"
	"fmt"
	"math/rand"
)

var src = "fmt.Println abfmt.Println abfmt.Println abfmt.Println abmt.Println abfmt.Println abfmt.Println abfmt.Println ab\nc"

func Hex32(n uint32) (ret [8]byte) {
	const board = "0123456789abcdef"
	shift := uint32(32)
	for i := 0; i < 8; i ++ {
		shift -= 4
		s := ((n >> shift) & 0xf)
		b := board[s]
		ret[i] = b
	}
	return
}

func Benchmark_sprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := uint32(rand.Int31())
		s := fmt.Sprintf("%08x %s\n", n, src)
		_ = s
	}
}

func Benchmark_sprintf_custom_hex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := uint32(rand.Int31())
		s := fmt.Sprintf("%8s %s\n", Hex32(n), src)
		_ = s
	}
}

func Benchmark_default_buffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		n := uint32(rand.Int31())

		buf.WriteString(fmt.Sprintf("%08x", n))

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}

func Benchmark_pre_alloc_buffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, len(src) + 32))
		n := uint32(rand.Int31())

		buf.WriteString(fmt.Sprintf("%08x", n))

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}

func Benchmark_pre_alloc_buffer_align(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 256))
		n := uint32(rand.Int31())

		buf.WriteString(fmt.Sprintf("%08x", n))

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}

func Benchmark_hex_sprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 256))
		n := uint32(rand.Int31())

		buf.WriteString(fmt.Sprintf("%08x", n))

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}

func Benchmark_hex_custom_make(b *testing.B) {
	// SLOW!
	var nn = 256
	for i := 0; i < b.N; i++ {
		s := make([]byte, 0, nn)
		_ = s
	}
}
func Benchmark_hex_custom_make2(b *testing.B) {
	// FAST!
	const nn = 256
	for i := 0; i < b.N; i++ {
		s := make([]byte, 0, nn)
		_ = s
	}
}

func Benchmark_hex_custom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf *bytes.Buffer
		if len(src) < 40 {
			buf = bytes.NewBuffer(make([]byte, 0, 64))
		} else if len(src) < 200 {
			buf = bytes.NewBuffer(make([]byte, 0, 256))
		} else if len(src) < 800 {
			buf = bytes.NewBuffer(make([]byte, 0, 1024))
		} else {
			nn := len(src) + len(src) / 16
			buf = bytes.NewBuffer(make([]byte, 0, nn))
		}

		n := uint32(rand.Int31())

		x := Hex32(n)
		buf.Write(x[:])

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}

func Benchmark_hex_custom2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, 512))
		n := uint32(rand.Int31())

		x := Hex32(n)
		buf.Write(x[:])

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}

func Benchmark_hex_custom3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(make([]byte, 0, len(src) + len(src)/8))
		n := uint32(rand.Int31())

		x := Hex32(n)
		buf.Write(x[:])

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}

func Benchmark_hex_custom_default_buffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		n := uint32(rand.Int31())

		x := Hex32(n)
		buf.Write(x[:])

		buf.WriteByte(' ')
		buf.WriteString(src)
		buf.WriteByte('\n')
		s := buf.Bytes()
		_ = s
	}
}
