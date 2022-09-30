
# build stage
FROM golang:1.18 as builder
WORKDIR /service/graphql_lambda
ADD . /service/graphql_lambda

RUN ls
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/main -mod vendor  ./cmd/graphql_lambda
#copying graphql schema defenition
COPY ./assets/schema.graphql ./build/assets/



# deployment image
FROM public.ecr.aws/lambda/provided:al2 as deploy

WORKDIR /service/
COPY --from=builder /service/graphql_lambda/build/* ./
RUN ls
ENTRYPOINT ["/main"]
