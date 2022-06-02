
# build stage
FROM golang:1.17 as builder
WORKDIR /service/userservice
COPY . /service/userservice
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/userservice -mod vendor  ./cmd/userservice

# deployment image
FROM alpine:3.10 as deploy
WORKDIR /service/
COPY --from=builder /service/userservice/build/* ./
CMD [ "sh", "-c", "sleep 5; /service/userservice" ]
