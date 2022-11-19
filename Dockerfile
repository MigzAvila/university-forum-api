# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /backend

# Download necessary Go modules
# COPY go.mod ./
# COPY go.sum ./


COPY . /backend
COPY . ./
RUN go mod download

WORKDIR /backend/cmd/api

RUN go build -o /api-backend

EXPOSE 4000

CMD [ "/api-backend" ]
