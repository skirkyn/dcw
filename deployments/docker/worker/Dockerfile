FROM golang:1.21.5

COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./
COPY main.go ./
COPY configs ./

ENV COMMON_CONFIG_LOC=./common.json
ENV CONTROLLER_CONFIG_LOC=./controller.json
ENV WORKER_CONFIG_LOC=./worker.json

RUN CGO_ENABLED=0 GOOS=linux go build -o /worker
RUN cp configs/* worker/

EXPOSE 50000

CMD ["/worker", "worker"]