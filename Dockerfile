FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main cmd/statbot/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY --from=builder /app/main /opt/statbot
COPY --from=builder /app/page /opt/page 
COPY --from=builder /app/.env /opt/.env
CMD ["/opt/statbot"]