all: local containers

local:
	GOOS="linux" GOARCH="amd64" go build -o bin/parrot-linux-amd64 .
	GOOS="linux" GOARCH="arm64" go build -o bin/parrot-linux-arm64 .
	GOOS="freebsd" GOARCH="amd64" go build -o bin/parrot-freebsd-amd64 .
	GOOS="freebsd" GOARCH="arm64" go build -o bin/parrot-freebsd-arm64 .
containers:
	podman build --jobs=2 --platform=linux/amd64,linux/arm64 --manifest parrot .
containers-publish:
	# you need to `podman login src.tty.cat` first
	podman manifest push localhost/parrot docker://src.tty.cat/home.arpa/parrot:latest
