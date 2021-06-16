FROM alpine:20210212@sha256:a81f67aca260e41b2260a4cf91ac4435dcad62ccc108ab19ae3077aafa1a4722

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
