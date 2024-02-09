FROM golang:alpine as builder

RUN go version
ENV GOPATH=/

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver .

FROM alpine as runner

WORKDIR /app

COPY --from=builder /build/httpserver .

CMD ["./httpserver"]