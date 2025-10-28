FROM alpine:3.22

ARG TARGETARCH

COPY --chmod=755 linux/${TARGETARCH}/openfeature /usr/local/bin/openfeature

ENTRYPOINT ["/usr/local/bin/openfeature"]
