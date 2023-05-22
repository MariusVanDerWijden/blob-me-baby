FROM golang:1.20.1-buster AS builder

# Override the default value of GOOS and GOARCH when building the Docker image using the --build-arg flag
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /build
# Copy and download dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o blob-me-baby main.go

FROM alpine:latest

COPY --from=builder /build/blob-me-baby /blob-me-baby

ENTRYPOINT ["/blob-me-baby"]