FROM golang:alpine as builder

WORKDIR /build 
COPY . .
RUN go build -o main .

# STAGE 2

FROM alpine

RUN adduser -D  appuser

WORKDIR /app 
COPY --from=builder /build/main /app/

USER appuser

EXPOSE 3000


CMD ["./main"]
