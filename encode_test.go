package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/bits"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

var tests = []struct {
	input    []byte
	want     []byte
	wordsize int
}{
	{
		[]byte{0xff, 0xff, 0xff},
		[]byte{0x3f, 0x3f, 0x3f, 0x3f},
		1,
	},
	{
		[]byte{0x00, 0x00, 0x00},
		[]byte{0x00, 0x00, 0x00, 0x00},
		1,
	},
	{
		[]byte{0xff, 0xff, 0xff},
		[]byte{0x3f, 0xff, 0x3f, 0xf0},
		2,
	},
	{
		[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		[]byte{0x3f, 0x3f, 0x3f, 0x3f, 0x3f, 0x3f, 0x3f, 0x3f, 0x3f, 0x30},
		1,
	},
	{
		[]byte{0xff, 0x00, 0xff},
		[]byte{0x3f, 0xc0, 0x0f, 0xf0},
		2,
	},
}

func TestPackTightlyString(t *testing.T) {
	for i, test := range tests {
		got := packTightly(test.input, test.wordsize)
		if !bytes.Equal(got, test.want) {
			t.Fatalf("test %v failed, want %b got %b", i, test.want, got)
		}
	}
}

func TestPackTightlyFast(t *testing.T) {
	for i, test := range tests {
		fmt.Printf("TEST: %v\n", i)
		got := packTightlyFast(test.input, test.wordsize)
		if countOnes(test.want) != countOnes(test.input) {
			panic("invalid test")
		}
		if !bytes.Equal(got, test.want) {
			t.Fatalf("test %v failed, want %b got %b", i, test.want, got)
		}
	}
}

func TestEncode(t *testing.T) {
	testEncode := func(data []byte) {
		if _, err := EncodeBlobs(data, false); err != nil {
			t.Fatal(err)
		}
		if _, err := EncodeBlobs(data, true); err != nil {
			t.Fatal(err)
		}
	}
	data := []byte{0xff, 0xff, 0xff}
	testEncode(data)

	rng := make([]byte, params.BlobTxBytesPerFieldElement*params.BlobTxFieldElementsPerBlob)
	rand.Read(rng)
	testEncode(rng)
}

// BenchmarkPackTight-8   	       3	 375693616 ns/op	2187486914 B/op	  135308 allocs/op
func BenchmarkPackTight(b *testing.B) {
	rng := make([]byte, params.BlobTxBytesPerFieldElement*params.BlobTxFieldElementsPerBlob)
	rand.Read(rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packTightly(rng, 32)
	}
}

// BenchmarkPackTightFast-8   	      48	  23178028 ns/op	 6056594 B/op	  131090 allocs/op
func BenchmarkPackTightFast(b *testing.B) {
	rng := make([]byte, params.BlobTxBytesPerFieldElement*params.BlobTxFieldElementsPerBlob)
	rand.Read(rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packTightlyFast(rng, 32)
	}
}

// BenchmarkPack-8   	   18536	     65399 ns/op	  303104 B/op	       2 allocs/op
func BenchmarkPack(b *testing.B) {
	rng := make([]byte, params.BlobTxBytesPerFieldElement*params.BlobTxFieldElementsPerBlob)
	rand.Read(rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pack(rng)
	}
}

func countOnes(data []byte) int {
	var res int
	for i := 0; i < len(data); i++ {
		res += bits.OnesCount(uint(data[i]))
	}
	return res
}
