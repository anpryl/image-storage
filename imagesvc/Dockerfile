FROM golang:1.10-alpine as builder
RUN apk update && apk add git && apk add binutils && apk add ca-certificates
RUN adduser -D -g imagesvc imagesvc
RUN mkdir -p /app/images && chown -R imagesvc /app && chmod 0700 /app && chmod 0700 /app/images
COPY . $GOPATH/src/github.com/image-storage/imagesvc
WORKDIR $GOPATH/src/github.com/image-storage/imagesvc
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-386 .
RUN chmod +x dep-linux-386 && \
    ./dep-linux-386 ensure && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app/imagesvc ./cmd/imagesvc/main.go && \
    cd /app && \
    strip --strip-unneeded imagesvc

FROM scratch
EXPOSE 4000
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=imagesvc:imagesvc /app /app
USER imagesvc
VOLUME /app/images
WORKDIR /app
ENTRYPOINT ["./imagesvc", "-folder", "/app/images"]
