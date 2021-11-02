FROM golang:latest as builder

ENV GO111MODULE=on

WORKDIR $GOPATH/src/github.com/GlassNode/be-challenge


COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build --tags apiserver ./main/api-server.go

ENTRYPOINT ["./api-server"]
