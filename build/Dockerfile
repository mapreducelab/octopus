FROM golang:1.11.1-alpine3.8 as builder

LABEL maintainer="antonputra@icloud.com"

RUN apk add --update git make curl && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/octopus

COPY . .

RUN make

FROM alpine:3.8

RUN addgroup -S octopus && \
    adduser -S -G octopus octopus

COPY --from=builder /go/src/octopus/octopus /home/octopus/

RUN chown -R octopus:octopus /home/octopus

USER octopus

ENTRYPOINT ["/home/octopus/octopus"]
