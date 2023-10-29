FROM golang:1.21

ARG PORT=8080
ARG COMMIT=unknown
ARG VERSION=0.0.1

EXPOSE ${PORT}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags "-X 'main.version=$VERSION' -X 'main.commit=$COMMIT'" -o /trapa-api

CMD ["/trapa-api"]