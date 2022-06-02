generate:
	# (rm -rf internal/proto)
	# protoc --go_out=plugins=grpc:. api/*.proto 
	(cd internal/configs && go generate && goimports -w *.go)
	(cd internal/gateway/configs && go generate && goimports -w *.go)


run:
	docker-compose up

configure:
	go install github.com/gokultp/go-envparser
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	grep -v "firstcontributions" /etc/hosts >> /tmp/hosts
	echo "172.30.1.6 api.firstcontributions.com" >> /tmp/hosts
	echo "172.30.1.8 explorer.firstcontributions.com" >> /tmp/hosts

	mv /tmp/hosts /etc/hosts