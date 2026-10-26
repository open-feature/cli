FROM alpine:3.24@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

ARG TARGETARCH

COPY --chmod=755 linux/${TARGETARCH}/openfeature /usr/local/bin/openfeature

ENTRYPOINT ["/usr/local/bin/openfeature"]
