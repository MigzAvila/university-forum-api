#  syntax=docker/dockerfile:1

# Fetch the base image
FROM golang:1.16-alpine 

# create a folder
WORKDIR /backend

# copy files 
COPY . /backend
COPY . ./
# install go dependencies
RUN go mod download

# build the project
RUN go build -o /api-backend ./cmd/api 

# expose the project to the public
EXPOSE 4000

# run the project
CMD [ "/api-backend", "--port=4000" ]
