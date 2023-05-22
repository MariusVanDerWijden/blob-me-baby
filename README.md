# Blob-me-baby

Blob-me-baby is a simple Go-based web server that provides an API endpoint for sending arbitrary data and receiving well formed
blobs, commitments and proofs.

## Installation

To install blob-me-baby, you need to have Go installed on your machine. Once Go is installed, you can clone the repository and build the project.

```bash

git clone https://github.com/MariusVanDerWijden/blob-me-baby.git
cd blob-me-baby
go build .
```

## Usage

Once you've built the project, you can start the server with:

```bash
./blob-me-baby
```

This will start the server on port 8080.

## Endpoints

Currently, blob-me-baby provides the following API endpoints: `/random` and `/encode/{data}`.

`/encode/{data}` is a GET endpoint that takes hexadecimal data as a path parameter and returns the computed commitments and proofs of the given data. 
The data is first transformed into Ethereum blobs, then commitments and proofs are computed, and finally returned as a JSON response.
