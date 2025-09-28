# Set target options
GOOS 		:= linux 
GOARCH 	:= arm 
GOARM 	:= 7

.PHONY: release debug help

release:
	@echo "Release mode ..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" sensors.go

debug:
	@echo "Debug mode..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -gcflags="all=-N -l" sensors.go

help:
	@echo "Available options:"
	@echo "  make           -> Compile in release mode"
	@echo "  make debug     -> Compile in debugen mode"
	@echo "  make help      -> Show this"
