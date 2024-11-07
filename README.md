# Fenrir
A modular game engine built in Go.

# Build proto generated code
```shell
cd internal/network && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protocol.proto
```

# Build the Fenrir server
```shell
go build -o build/fenrir_server.exe cmd/server/main.go
```

# Build the Fenrir client
```shell
go build -o build/fenrir_server.exe cmd/server/main.go
```