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

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder keruu /keruu
ENTRYPOINT ["/keruu"]
