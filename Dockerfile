FROM ubuntu:24.04 AS builder
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        binutils \
        upx \
    && rm -rf /var/lib/apt/lists/*
COPY keruu .
RUN strip keruu && upx -q -9 keruu
RUN echo "nobody:x:10101:10101:Nobody:/:" > /etc_passwd

FROM gcr.io/distroless/base-debian11
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc_passwd /etc/passwd
COPY --from=builder --chown=10101:0 keruu /keruu
USER nobody
ENTRYPOINT ["/keruu"]
