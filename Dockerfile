FROM golang:1.19-alpine AS build

WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .

ARG SERVICE
ENV EVENTSTORE_HOST=host.minikube.internal
ENV MONGODB_HOST=host.minikube.internal
ENV MONGODB_USER=root
ENV MONGODB_PASS=example

RUN go build -o /server ./${SERVICE}/cmd/.

ENTRYPOINT ["/server"]