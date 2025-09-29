FROM golang:1.25-alpine AS builder

WORKDIR /build 

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build main.go diberi nama server
RUN go build -o server ./cmd/main.go


FROM alpine:3.22

# copy hasil build (server) dari stage builder ke directory /app/server
WORKDIR /app
COPY --from=builder /build/server ./server

RUN chmod +x server

EXPOSE 2409

CMD [ "/app/server" ]