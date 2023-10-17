FROM golang:1.21

ARG PORT=8080

EXPOSE ${PORT}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /trapa-api

CMD ["/trapa-api"]