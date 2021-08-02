FROM golang:1.16-alpine3.14 AS builder

WORKDIR /src
COPY . .
RUN go build -ldflags="-s -w" -o ./kube-image-clone-controller *.go

FROM alpine:3.14.0

WORKDIR /opt
COPY --from=builder /src/kube-image-clone-controller /opt

ENTRYPOINT ["sh","-c","./kube-image-clone-controller"]



