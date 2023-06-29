FROM gcr.io/distroless/base-debian11
COPY keruu /usr/local/bin/keruu
ENTRYPOINT ["keruu"]
