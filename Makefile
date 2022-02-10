generate:
	# (rm -rf internal/proto)
	# protoc --go_out=plugins=grpc:. api/*.proto 
	(cd internal/configs && go generate && goimports -w *.go)
	(cd internal/gateway/configs && go generate && goimports -w *.go)


run:
	docker-compose up

configure:
	go install github.com/gokultp/go-envparser
	go install github.com/golang/protobuf/proto
	go install github.com/golang/protobuf/protoc-gen-go
	grep -v "172.30.1." /etc/hosts >> /tmp/hosts
	echo "172.30.1.6 api.firstcontributions.com" >> /tmp/hosts
	echo "172.30.1.8 explorer.firstcontributions.com" >> /tmp/hosts

	mv /tmp/hosts /etc/hosts