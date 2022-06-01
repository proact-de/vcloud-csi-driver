FROM alpine:20220328@sha256:c8420124adc6f31031a88c85f53e1e78f419b1ac46b827fd1d4d0fdae2c7e3aa

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
