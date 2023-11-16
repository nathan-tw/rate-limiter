FROM golang:alpine

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o app .

ENTRYPOINT [ "./app" ]