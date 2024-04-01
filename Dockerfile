FROM golang:1.22.1-alpine
WORKDIR /app


COPY . .
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /starter-go-app

EXPOSE 8080
CMD ["/starter-go-app"]