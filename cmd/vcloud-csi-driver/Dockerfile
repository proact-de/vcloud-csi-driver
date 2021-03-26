FROM alpine:20210212@sha256:fa3bd1cb8b0d2d6a4de1ea7e52dffee36896bc2b1566e9a89c16637051467225

RUN apk add --no-cache \
    e2fsprogs \
    e2fsprogs-extra \
    xfsprogs \
    xfsprogs-extra \
    blkid \
    findmnt \
    mailcap \
    ca-certificates

EXPOSE 8080
CMD ["/usr/bin/vcloud-csi-driver"]

COPY bin/vcloud-csi-driver /usr/bin/vcloud-csi-driver
