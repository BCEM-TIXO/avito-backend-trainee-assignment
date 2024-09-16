FROM golang:1.22-alpine AS build

WORKDIR /src

COPY ./src .

RUN go mod download

WORKDIR /src/cmd/app

RUN go build -o main .

FROM alpine:latest

WORKDIR /root/

EXPOSE 8080


COPY --from=build /src/cmd/app/main .
COPY --from=build /src/migrations /root/migrations

ENTRYPOINT ["./main"]