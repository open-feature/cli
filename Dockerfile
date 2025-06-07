FROM alpine:3.22

COPY ./openfeature usr/local/bin/openfeature

RUN chmod +x /usr/local/bin/openfeature

ENTRYPOINT ["/usr/local/bin/openfeature"]
