FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install Goose for database migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go build -v -o /usr/local/bin/app .

ENTRYPOINT ["bash", "-c", "goose -dir ./db postgres $DATABASE_URL up && app"]
