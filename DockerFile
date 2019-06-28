
FROM golang:latest as builder

LABEL version="v0.1.0"
LABEL repository="https://github.com/w9jds/beholder-bot"
LABEL homepage="https://github.com/w9jds/beholder-bot"
LABEL maintainer="Jeremy Shore <w9jds@github.com>"

WORKDIR /go/src/beholder/cmd

COPY . /go/src/beholder

RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o beholder
RUN curl -o ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt

FROM scratch

WORKDIR /go/src/beholder/cmd

COPY --from=builder /go/src/beholder/cmd/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/beholder/cmd/beholder /go/src/beholder/cmd

CMD ["./beholder"]