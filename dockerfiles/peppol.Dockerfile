# BASE IMAGE
FROM golang:1.22-alpine as base

WORKDIR /app

COPY . /app
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/peppol

# RUNNER
FROM alpine:latest
WORKDIR /app
COPY --from=base /app/peppol .

ENTRYPOINT ["./peppol"]

