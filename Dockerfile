FROM golang:1.18.8-alpine3.16 AS BuildStage

WORKDIR /crud
COPY . .
RUN go mod download
EXPOSE 8000
RUN go build -o /test main.go

FROM alpine:latest

WORKDIR /
COPY --from=BuildStage /test /test
EXPOSE 8000
ENTRYPOINT ["/test"]


