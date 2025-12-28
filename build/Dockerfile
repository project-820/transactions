FROM golang:1.24 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/transactions ./cmd
RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

FROM alpine:3.21 AS final

WORKDIR /

RUN apk add --no-cache ca-certificates

COPY --from=build /bin/transactions /transactions
COPY --from=build /go/bin/grpc-health-probe /usr/local/bin/grpc-health-probe

EXPOSE 8082

ENTRYPOINT ["/transactions"]