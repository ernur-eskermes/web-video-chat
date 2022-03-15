FROM golang:1.17-alpine AS builder

RUN go version

COPY . /github.com/ernur-eskermes/web-video-chat
WORKDIR /github.com/ernur-eskermes/web-video-chat

RUN go mod download
RUN GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/ernur-eskermes/web-video-chat/.bin/app .
COPY --from=0 /github.com/ernur-eskermes/web-video-chat/configs configs/
COPY --from=0 /github.com/ernur-eskermes/web-video-chat/templates templates/

CMD ["./app"]