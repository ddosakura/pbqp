.PHONY: clean
clean: proto
	rm ./*.sock

.PHONY: client
client: proto
	go run server/main.go

.PHONY: server
server: proto
	go run server/main.go

.PHONY: proto
proto:
	protoc --go_out=. proto/msg/msg.proto
