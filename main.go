package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting server at port 80")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/random", handleRandom).
		Methods("GET")
	router.HandleFunc("/encode/{data}", handleEncode).
		Methods("GET")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}

func handleRandom(rw http.ResponseWriter, req *http.Request) {
	// TODO impl
}

// CommitmentAndProof is the response to the encode endpoint
type CommitmentAndProof struct {
	Commitments     []types.KZGCommitment
	VersionedHashes []common.Hash
	AggregatedProof []types.KZGProof
}

func handleEncode(rw http.ResponseWriter, req *http.Request) {
	data := common.FromHex(mux.Vars(req)["data"])
	blobs := encodeBlobs(data)
	commits, versionedHashes, aggProof, err := blobs.ComputeCommitmentsAndProofs()
	if err != nil {
		fmt.Printf("Error computing commitment: %v\n", err)
		rw.WriteHeader(500)
		return
	}

	result := CommitmentAndProof{
		Commitments:     commits,
		VersionedHashes: versionedHashes,
		AggregatedProof: aggProof,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("Error marshalling result: %v\n", err)
		rw.WriteHeader(500)
		return
	}
	rw.Write(resp)
}

func encodeBlobs(data []byte) types.Blobs {
	blobs := []types.Blob{{}}
	blobIndex := 0
	fieldIndex := -1
	for i := 0; i < len(data); i += 31 {
		fieldIndex++
		if fieldIndex == params.FieldElementsPerBlob {
			blobs = append(blobs, types.Blob{})
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
