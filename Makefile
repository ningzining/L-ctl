TARGET_NAME = L-ctl

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(TARGET_NAME) .

.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(TARGET_NAME) .

.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(TARGET_NAME).exe .

.PHONY: build-clean
build-clean:
	@rm -f $(CURDIR)/$(TARGET_NAME)
	@rm -f $(CURDIR)/$(TARGET_NAME).exe