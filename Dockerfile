# build stage
FROM golang:1.18.10-alpine3.17 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .

# only specify what a port will intend to be used, don't force it
EXPOSE 8080
CMD [ "/app/main" ]
