FROM golang:latest

WORKDIR /go/src
RUN ln -sf /bin/bash /bin/sh

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

EXPOSE 8080

CMD ["tail", "-f", "/dev/null"]



