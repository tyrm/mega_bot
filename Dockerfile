FROM golang:1.14 AS builder
RUN go get github.com/markbates/pkger/cmd/pkger

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN pkger && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mega_bot

#FROM alpine:latest
FROM scratch
COPY --from=builder /etc/ssl /etc/ssl

COPY --from=builder /app/mega_bot /mega_bot
CMD ["/mega_bot"]