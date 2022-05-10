FROM golang:1.18-alpine AS builder
RUN apk update && apk add git build-base
WORKDIR /go/src/github.com/natemago/card-games-api

ADD app ./app
ADD cmd ./cmd
ADD config ./config
ADD repositories ./repositories
ADD rest ./rest
ADD errors ./errors
ADD go.mod ./
ADD go.sum ./
ADD main.go ./

RUN go build


FROM alpine:3

EXPOSE 8080

WORKDIR /root

COPY --from=builder /go/src/github.com/natemago/card-games-api/card-games-api ./card-games-api
RUN chmod +x ./card-games-api

CMD ["./card-games-api"]