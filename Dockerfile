FROM golang:1.18-alpine AS builder
RUN apk update && apk add git
WORKDIR /go/src/toggl.com/services/card-games-api

ADD app ./app
ADD cmd ./cmd
ADD config ./config
ADD repositories ./repositories
ADD rest ./rest
ADD go.mod ./
ADD go.sum ./
ADD main.go ./

RUN go build
RUN ls -la
RUN pwd


FROM alpine:3

EXPOSE 8080

WORKDIR /root

COPY --from=builder /go/src/toggl.com/services/card-games-api/card-games-api ./card-games-api
RUN chmod +x ./card-games-api
RUN ls -la
RUN pwd

CMD ["./card-games-api"]