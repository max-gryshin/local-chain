# builder image
FROM golang:1.16.3 AS builder
WORKDIR /app
COPY /backend /app
RUN go mod download && go get -u -v -f all
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o local-chain

FROM alpine:3.11.3
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/local-chain .
EXPOSE 8080

ENTRYPOINT [ "/local-chain" ]
