# builder image
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY . .

# RUN go mod download && go get -u -v -f all
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -installsuffix cgo -o local-chain ./cmd/local-chain

# Create directories for Raft data and BoltDB
RUN mkdir -p /db

FROM alpine:3.11.3
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/local-chain .
EXPOSE 8001 9001

ENTRYPOINT [ "/local-chain" ]
