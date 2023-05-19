# shared variables
GIT_SHA := $(shell git rev-parse HEAD)
GOMODCACHE := $(shell go env GOMODCACHE)
GOCACHE := $(shell go env GOCACHE)
GOOS 	:= linux
GOARCH  := amd64

REPO	:= zeusfyi
NAME    := blobme
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest

docker.pubbuildx:
	@ docker buildx build \
 	-t ${IMG} \
  	-t ${LATEST} \
  	--build-arg GOMODCACHE=${GOMODCACHE} \
  	--build-arg GOCACHE=${GOCACHE} \
  	--build-arg GOOS=${GOOS} \
  	--build-arg GOARCH=${GOARCH} \
  	--platform=${GOOS}/${GOARCH} \
  	-f ./Dockerfile . --push

docker.debug:
	docker run -it --entrypoint /bin/bash zeusfyi/blobme:latest