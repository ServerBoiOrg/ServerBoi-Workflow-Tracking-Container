FROM golang:1.16-alpine

LABEL org.opencontainers.image.source="https://github.com/serverboiorg/status-monitor"

WORKDIR /status

COPY status/go.mod ./
COPY status/go.sum ./

RUN go mod download

COPY status/*.go ./

RUN export CGO_ENABLED=0
RUN go build -o /serverboi-status-monitor

EXPOSE 7032/tcp

CMD [ "/serverboi-status-monitor" ]