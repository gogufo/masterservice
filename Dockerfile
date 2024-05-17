FROM golang:1.21.6 AS builder

WORKDIR /go/bin

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .


RUN go build -o /go/bin/masterservice  *.go


FROM ubuntu

COPY --from=builder /go/bin/masterservice /go/bin/masterservice
COPY --from=builder /go/bin/config/ /go/bin/config/


WORKDIR /go/bin/

EXPOSE 5300

ENTRYPOINT ["/go/bin/masterservice"]
