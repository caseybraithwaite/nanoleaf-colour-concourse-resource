FROM golang:1.20.0-alpine as builder

ENV GO111MODULE=on

WORKDIR /build

COPY src/ .

RUN go build -o assets/out out/out.go

FROM alpine:latest

COPY --from=builder /build/assets/out /opt/resource/out
COPY src/in/in.sh /opt/resource/in
COPY src/check/check.sh /opt/resource/check

RUN chmod +x -R /opt/resource