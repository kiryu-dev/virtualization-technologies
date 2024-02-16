FROM golang:alpine as builder

RUN go version
ENV GOPATH=/

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver .

FROM alpine as runner

ARG USERNAME="kira"
RUN adduser --disabled-password $USERNAME

WORKDIR /app

USER $USERNAME

COPY --from=builder --chown=kira /build/httpserver .

RUN echo "hi, i am" `whoami`

CMD ["./httpserver"]