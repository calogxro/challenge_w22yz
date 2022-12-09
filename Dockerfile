FROM golang:1.19-alpine AS build
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -o /server .

EXPOSE 8080

ENTRYPOINT ["/server"]