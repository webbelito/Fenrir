# Fenrir
A modular game engine built in Go.

# Now requires gcc to build

- Visit https://code.visualstudio.com/docs/cpp/config-mingw#_prerequisites and go to the "direct link to the installer"
- Install the application, you'll get a bash window to open once the installer completes.
- Run the pacman command: 

```shell 
pacman -S --needed base-devel mingw-w64-ucrt-x86_64-toolchain
``` 

- Press Enter to use default=all
- Press Y to accept
- Set the PATH in Windows
- Verify the gcc version

```shell
gcc --version
```

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