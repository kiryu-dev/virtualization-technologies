FROM golang:alpine

RUN go version
ENV GOPATH=/

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver .

CMD ["./httpserver"]