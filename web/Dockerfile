FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./web

FROM gcr.io/distroless/base-debian11:nonroot
WORKDIR /

COPY --from=builder /app ./
COPY --from=builder /web/views ./views

EXPOSE 8080

CMD ["./app"]