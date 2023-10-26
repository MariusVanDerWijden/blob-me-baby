package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
)

func EncodeBlobs(data []byte, compressed bool) (*CommitmentAndProof, error) {
	var packed []byte
	if compressed {
		packed = packTightly(data, 32)
	} else {
		packed = pack(data)
	}
	blobs := encodeBlobs(packed)

	var result CommitmentAndProof
	for _, blob := range blobs {
		commit, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			return nil, err
		}
		proof, err := kzg4844.ComputeBlobProof(blob, commit)
		if err != nil {
			return nil, err
		}

		result.Commitments = append(result.Commitments, commit)
		result.Blobs = append(result.Blobs, blob)
		result.AggregatedProof = append(result.AggregatedProof, proof)
		result.VersionedHashes = append(result.VersionedHashes, kZGToVersionedHash(commit))
	}
	return &result, nil
}

// encodeBlobs expects the data to be packed correctly.
func encodeBlobs(data []byte) []kzg4844.Blob {
	// Put the packed data into blobs
	var (
		blobs      = []kzg4844.Blob{{}}
		blobIndex  = 0
		fieldIndex = -1
	)
	for i := 0; i < len(data); i += 32 {
		fieldIndex++
		if fieldIndex == params.BlobTxFieldElementsPerBlob {
			blobs = append(blobs, kzg4844.Blob{})
			blobIndex++
			fieldIndex = 0
		}
		max := i + 32
		if max > len(data) {
			max = len(data)
		}
		copy(blobs[blobIndex][fieldIndex*32:], data[i:max])
	}
	return blobs
}

func pack(data []byte) []byte {
	result := make([]byte, 0, len(data))
	for i := 0; i < len(data); i += 31 {
		max := i + 31
		if max > len(data) {
			max = len(data)
		}
		result = append(result, 0x00)
		result = append(result, data[i:max]...)
	}
	return result
}

func packTightly(data []byte, wordsize int) []byte {
	s := ""
	for i := 0; i < len(data); i++ {
		s += fmt.Sprintf("%08b", data[i])
	}
	for i := 0; i < len(s); i += 8 * wordsize {
		s = s[:i] + "00" + s[i:]
	}
	if missing := len(s) % (8 * wordsize); missing != 0 {
		s = s + strings.Repeat("0", (8*wordsize)-missing)
	}
	res := make([]byte, 0, len(s)/8)
	for i := 0; i < len(s); i += 8 {
		b, err := strconv.ParseUint(s[i:i+8], 2, 8)
		if err != nil {
			panic("conversion failed")
		}
		res = append(res, byte(b))
	}
	return res
}

var blobCommitmentVersionKZG uint8 = 0x01

// kZGToVersionedHash implements kzg_to_versioned_hash from EIP-4844
func kZGToVersionedHash(kzg kzg4844.Commitment) common.Hash {
	h := sha256.Sum256(kzg[:])
	h[0] = blobCommitmentVersionKZG

	return h
}
