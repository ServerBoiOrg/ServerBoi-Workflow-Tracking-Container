FROM golang:1.16-alpine

LABEL maintainer="awlsring"

WORKDIR /workflow-tracking

COPY workflow-tracking/go.mod ./
COPY workflow-tracking/go.sum ./

RUN go mod download

COPY workflow-tracking/*.go ./

RUN go build -o /serverboi-workflow-tracking

CMD [ "/serverboi-workflow-tracking" ]