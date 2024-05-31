BINARY := backend

$(BINARY): *.go go.mod go.sum
	go build -ldflags="-s -w" -trimpath -o $@
	-upx $@
