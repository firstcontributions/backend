
# build stage
FROM golang:1.14 as builder
WORKDIR /service/profile
COPY . /service/profile
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/profile -mod vendor  ./cmd/profile

# deployment image
FROM alpine:3.10 as deploy
WORKDIR /service/
COPY --from=builder /service/profile/build/* ./
CMD [ "sh", "-c", "sleep 5; /service/profile" ]
