FROM golang:1.21.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY /cmd/ ./cmd
COPY configs/ ./configs
RUN mv ./configs/* ./
ENV COMMON_CONFIG_LOC=common.json
ENV CONTROLLER_CONFIG_LOC=controller.json
ENV WORKER_CONFIG_LOC=worker.json

RUN apk update
RUN apk add libzmq zeromq-dev zeromq gcc build-base pkgconfig
RUN CGO_ENABLED=1  GOOS=linux go build -o worker ./cmd/worker

RUN chmod +x ./worker

CMD ["./worker"]