FROM amd64/golang:1.17-alpine

LABEL org.opencontainers.image.source="https://github.com/ServerBoiOrg/ServerBoi-Workflow-Tracking-Container"

RUN apk update
RUN apk add git

WORKDIR /workflow

COPY workflow-tracking/go.mod ./
COPY workflow-tracking/go.sum ./

RUN git clone https://github.com/ServerBoiOrg/ServerBoi-Lambdas-Go ./ServerBoi-Lambdas-Go

RUN go mod download

COPY workflow-tracking/*.go ./

RUN go build -o /serverboi-workflow-tracking

CMD [ "/serverboi-workflow-tracking" ]
