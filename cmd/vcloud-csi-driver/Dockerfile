FROM alpine:20210730@sha256:12341c2503d96f13a17b0c935ae9090034e317afdfc934f1e0df1a5be813e73e

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
