FROM golang:1.10-alpine as builder
RUN apk update && apk add git && apk add binutils && apk add ca-certificates
RUN adduser -D -g '' authsvc
COPY . $GOPATH/src/github.com/image-storage/authsvc
WORKDIR $GOPATH/src/github.com/image-storage/authsvc
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-386 .
RUN chmod +x dep-linux-386 && \
    ./dep-linux-386 ensure && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o $GOPATH/bin/authsvc . && \
    cd $GOPATH/bin && \
    strip --strip-unneeded authsvc

FROM scratch
EXPOSE 4000
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/authsvc /go/bin/authsvc
USER authsvc
ENTRYPOINT ["/go/bin/authsvc"]
