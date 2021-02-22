FROM alpine:20210212@sha256:78b6d0a9e7a5b3d2aa1c10932f08272d360a9b35e80b3cc5fb04b30e04fecf4b

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
