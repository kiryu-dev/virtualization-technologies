FROM golang:alpine as builder

RUN go version
ENV GOPATH=/

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver ./cmd/app/main.go

FROM alpine as runner

ARG USERNAME="kira"
RUN adduser --disabled-password $USERNAME

WORKDIR /app

USER $USERNAME

ARG SERVER_CONFIG_PATH

COPY --from=builder --chown=kira /build/httpserver /build/${SERVER_CONFIG_PATH} ./

RUN echo "hi, i am" `whoami`

CMD ["./httpserver"]