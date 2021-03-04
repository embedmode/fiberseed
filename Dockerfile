FROM golang:1.16-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main /
COPY --from=builder /build/.env /

# Command to run
ENTRYPOINT ["/main"]