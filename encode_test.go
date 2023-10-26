package main

import (
	"bytes"
	"crypto/rand"
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

/*
func TestPackTightlyFast(t *testing.T) {
	for i, test := range tests {
		got := packTightlyFast(test.input, test.wordsize)
		if countOnes(test.want) != countOnes(test.input) {
			panic("invalid test")
		}
		if !bytes.Equal(got, test.want) {
			t.Fatalf("test %v failed, want %b got %b", i, test.want, got)
		}

	}
}
*/

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

func BenchmarkPack(b *testing.B) {
	rng := make([]byte, params.BlobTxBytesPerFieldElement*params.BlobTxFieldElementsPerBlob)
	rand.Read(rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packTightly(rng, 32)
	}
}

func countOnes(data []byte) int {
	var res int
	for i := 0; i < len(data); i++ {
		res += bits.OnesCount(uint(data[i]))
	}
	return res
}
