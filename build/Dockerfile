#Builder
FROM golang:1.21.0-alpine AS builder

COPY . /github.com/RIP/
WORKDIR /github.com/RIP/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o server ./app/cmd/main.go

FROM scratch AS run

WORKDIR /docker-rip/

COPY --from=builder /github.com/RIP/server .

EXPOSE 80 433 8080

ENTRYPOINT ["./server"]