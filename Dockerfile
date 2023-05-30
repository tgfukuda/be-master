# build stage
FROM golang:1.18.10-alpine3.17 AS builder

WORKDIR /app
COPY . .
# alpine doesn't have a curl by default
RUN apk add curl
# get sqlc
RUN curl -L https://github.com/kyleconroy/sqlc/releases/download/v1.18.0/sqlc_1.18.0_linux_amd64.tar.gz | tar xvz
RUN /app/sqlc generate
RUN go build -o main main.go

# Run stage
FROM alpine:3.17
WORKDIR /app

COPY --from=builder /app/main .

COPY app.env .

COPY start.sh .
COPY wait-for.sh .

COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
