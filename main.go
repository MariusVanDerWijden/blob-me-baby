package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gorilla/mux"
)

var listeningPort = flag.String("port", "8080", "Port to listen on")

func main() {
	flag.Parse()
	fmt.Printf("Starting server at port %s\n", *listeningPort)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/random", handleRandom).
		Methods("GET")
	router.HandleFunc("/encode/{data}", handleEncode).
		Methods("GET")

	// Format the string with the port number
	portStr := fmt.Sprintf(":%v", *listeningPort)

	if err := http.ListenAndServe(portStr, router); err != nil {
		log.Fatal(err)
	}
}

func handleRandom(rw http.ResponseWriter, req *http.Request) {
	// TODO impl
}

// CommitmentAndProof is the response to the encode endpoint
type CommitmentAndProof struct {
	VersionedHashes []common.Hash
	Commitments     []kzg4844.Commitment
	Blobs           []kzg4844.Blob
	AggregatedProof []kzg4844.Proof
}

func handleEncode(rw http.ResponseWriter, req *http.Request) {
	data := common.FromHex(mux.Vars(req)["data"])

	result, err := EncodeBlobs(data)
	if err != nil {
		fmt.Printf("Error encoding blobs: %v\n", err)
		rw.WriteHeader(500)
		return
	}

	resp, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Error marshalling result: %v\n", err)
		rw.WriteHeader(500)
		return
	}
	_, err = rw.Write(resp)
	if err != nil {
		log.Fatalf("Error writing response: %v\n", err)
	}
}

func EncodeBlobs(data []byte) (*CommitmentAndProof, error) {
	blobs := encodeBlobs(data)

	var result CommitmentAndProof
	for _, blob := range blobs {
		result.Blobs = append(result.Blobs, blob)

		commit, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			return nil, err
		}
		result.Commitments = append(result.Commitments, commit)

		proof, err := kzg4844.ComputeBlobProof(blob, commit)
		if err != nil {
			return nil, err
		}
		result.AggregatedProof = append(result.AggregatedProof, proof)

		result.VersionedHashes = append(result.VersionedHashes, kZGToVersionedHash(commit))
	}
	return &result, nil
}

func encodeBlobs(data []byte) []kzg4844.Blob {
	blobs := []kzg4844.Blob{{}}
	blobIndex := 0
	fieldIndex := -1
	for i := 0; i < len(data); i += 31 {
		fieldIndex++
		if fieldIndex == params.BlobTxFieldElementsPerBlob {
			blobs = append(blobs, kzg4844.Blob{})
			blobIndex++
			fieldIndex = 0
		}
		max := i + 31
		if max > len(data) {
			max = len(data)
		}
		copy(blobs[blobIndex][fieldIndex*32:], data[i:max])
	}
	return blobs
}

var blobCommitmentVersionKZG uint8 = 0x01

// kZGToVersionedHash implements kzg_to_versioned_hash from EIP-4844
func kZGToVersionedHash(kzg kzg4844.Commitment) common.Hash {
	h := sha256.Sum256(kzg[:])
	h[0] = blobCommitmentVersionKZG

	return h
}
