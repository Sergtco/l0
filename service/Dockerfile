FROM golang:1.22

WORKDIR /app


COPY . .
RUN go mod download && go mod verify
RUN go build ./cmd/main.go
EXPOSE 6969
CMD ["./main"]
