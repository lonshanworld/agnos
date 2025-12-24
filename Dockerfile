FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src

# copy modules and download
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/server ./server
EXPOSE 8080
CMD ["/app/server"]
