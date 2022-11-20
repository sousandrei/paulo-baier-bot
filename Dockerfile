FROM golang:1.19.3-bullseye

WORKDIR /app

RUN apt install git
RUN go install github.com/go-task/task/v3/cmd/task@latest

COPY .git .git
COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .
COPY Taskfile.yaml .

RUN task build

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

ENV TELEGRAM_TOKEN=

COPY --from=0 /app/build/pbb .

ENTRYPOINT ["/app/pbb"]