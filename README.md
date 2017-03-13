# Drone-in-Drone Plugin

This plugin allows you to run `drone exec` inside of Drone by utilizing
Docker-in-Docker.

## Build

Build the binary with the following commands:

```
drone exec
```

[Drone CLI](http://readme.drone.io/usage/getting-started-cli/) is required

## Docker

Build the docker image with the following commands:

```
drone exec
docker build -t jmccann/drone-in-drone .
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
