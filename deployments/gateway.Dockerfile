
# build stage
FROM golang:1.14 as builder
WORKDIR /service/gateway
COPY . /service/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/gateway -mod vendor  ./cmd/gateway
#copying graphql schema defenition
COPY ./internal/gateway/graphql/schema.graphql ./build/



# deployment image
FROM alpine:3.10 as deploy
WORKDIR /service/
COPY --from=builder /service/gateway/build/* ./
CMD [ "sh", "-c", "sleep 5; /service/gateway" ]
