FROM golang:1.12-alpine3.9

ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN apk add --no-cache make git

WORKDIR /go/src/github.com/martinohmann/kubernetes-cluster-manager

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN make install

FROM scratch

COPY --from=0 /go/bin/kcm /kcm

ENTRYPOINT ["/kcm"]
