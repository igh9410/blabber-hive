FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application
COPY . .

RUN go build -o blabber-hive ./cmd/main.go

EXPOSE 8080

# Specify the command to run on container start
CMD ["/app/blabber-hive"]