FROM docker:1.13-dind

RUN curl http://downloads.drone.io/release/linux/amd64/drone.tar.gz | tar xz -C /bin
ADD drone-in-drone /bin/
ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/bin/drone-in-drone"]
