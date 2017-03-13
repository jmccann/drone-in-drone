# Drone-in-Drone Plugin

This plugin allows you to run `drone exec` inside of Drone by utilizing
Docker-in-Docker.

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t jmccann/drone-publish-plugin .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-docker' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e DRONE_COMMIT_SHA=d8dbe4d94f15fe89232e0402c6e8a0ddf21af3ab \
  -e PLUGIN_CONTEXT=path/to/code \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  --privileged \
  jmccann/drone-in-drone
```
