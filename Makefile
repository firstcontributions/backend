generate:
	(rm -rf internal/proto)
	protoc -I=api/ --go_out=plugins=grpc:. api/*.proto
	(cd internal/profile/configs && go generate && goimports -w *.go)
	(cd internal/configs && go generate && goimports -w *.go)
	(cd internal/gateway/configs && go generate && goimports -w *.go)


run:
	docker-compose up