package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
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
	router.HandleFunc("/tight/{data}", handleCompressed).
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
	result, err := EncodeBlobs(data, false)
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

func handleCompressed(rw http.ResponseWriter, req *http.Request) {
	data := common.FromHex(mux.Vars(req)["data"])
	result, err := EncodeBlobs(data, true)
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
