# builder
FROM golang:1.20.3-bullseye as builder
WORKDIR /work
ENV GOPROXY=direct
COPY go.mod go.sum /work/
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o app

# prod
FROM debian:bullseye-slim as prod
WORKDIR /app
RUN apt-get update && \
    apt-get install -y --force-yes --no-install-recommends apt-transport-https curl ca-certificates && \
    apt-get clean && \
    apt-get autoremove && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /work/app .
CMD ["./app"]
