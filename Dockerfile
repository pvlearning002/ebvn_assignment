
FROM golang:1.27-alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./command/api/main.go

FROM alpine:3.23.2 AS runtime

WORKDIR /app

COPY --from=build /opt/app/main ./main

EXPOSE 8080

CMD ["./main"]
