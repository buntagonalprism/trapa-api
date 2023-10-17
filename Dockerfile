FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /trapa-api

EXPOSE 3000

CMD ["/trapa-api"]