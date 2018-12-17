FROM golang:1.11-alpine as builder
RUN apk --no-cache add git

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /vgo/
COPY . .

RUN go get ./...
RUN go test ./...
RUN go build -o /bin/kubewire .

# Build runtime
FROM alpine:3.8 as runtime
MAINTAINER OpenSource PF <opensource@postfinance.ch>

COPY --from=builder /bin/kubewire /bin/kubewire

# Run as nobody:x:65534:65534:nobody:/:/sbin/nologin
USER 65534

CMD ["/bin/kubewire"]
