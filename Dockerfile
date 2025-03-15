# builder image
FROM golang:1.21.0-alpine as builder
WORKDIR /build
COPY . /build
RUN CGO_ENABLED=0 GOOS=linux go build -o authorization-service .

# generate clean, final image for end users
FROM alpine:3.18 as hoster

RUN apk add --no-cache tzdata
ENV TZ=Asia/Qyzylorda

COPY --from=builder /build/authorization-service ./authorization-service
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/migrations/ ./migrations/

# executable
ENTRYPOINT [ "./authorization-service" ]
