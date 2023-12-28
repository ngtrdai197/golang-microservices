FROM golang:1.21 as builder
WORKDIR /codebase

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

# Building app
RUN GOOS=linux CGO_ENABLED=0 go build -o codebase.app .

FROM alpine:3.18

RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/UTC /etc/localtime && \
    apk del tzdata

WORKDIR /app
EXPOSE 8080
EXPOSE 8090

RUN mkdir -p config
COPY --from=builder /codebase/config/.env /app/config/.
COPY --from=builder /codebase/codebase.app .