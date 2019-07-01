.PHONY: clean
clean:
	rm ./*.sock

.PHONY: client
client: proto
	go run client/main.go ./swap.sock

.PHONY: server
server: proto
	go run server/main.go ./swap.sock

.PHONY: proto
proto:
	protoc --go_out=. proto/msg/msg.proto
