# Build package by default with all general tools and things
all:
	make build

# Build minimal package only
build:
	go build

# Build debug version with more logs
debug:
	go build -tags=debug

# Format the source code
fmt:
	go fmt