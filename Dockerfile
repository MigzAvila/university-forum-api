#  syntax=docker/dockerfile:1
FROM golang:1.16-alpine 

WORKDIR /backend

COPY . /backend
COPY . ./
RUN go mod download

# WORKDIR /backend/cmd/api

RUN go build -o /api-backend ./cmd/api 

EXPOSE 4000

CMD [ "/api-backend", "--port=4000"]
