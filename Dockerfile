FROM golang:1.22.2-alpine3.19 as build

WORKDIR /app
COPY . /app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build --ldflags '-extldflags=-static' -o changelog-json main.go

FROM alpine:3.19.1
COPY --from=build /app/changelog-json /app/
ENTRYPOINT ["/app/changelog-json"]
