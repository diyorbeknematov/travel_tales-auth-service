FROM golang:1.22.2 AS builder

WORKDIR /travel-auth

COPY . .
RUN go mod download

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../travel_app .

FROM alpine:latest

WORKDIR /travel-auth 

COPY --from=builder /travel-auth/travel_app .
COPY --from=builder /travel-auth/logs/app.log ./logs/
COPY --from=builder /travel-auth/.env .

EXPOSE 8081

CMD [ "./travel_app" ]