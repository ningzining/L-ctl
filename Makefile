TARGET_NAME = L-ctl

.PHONY: build
build:
	@go build -o L-ctl.exe main.go

.PHONY: build.linux
build.linux:
	@SET GOOS=linux
	@SET GOARCH=amd64
	@go build -o L-ctl main.go

.PHONY: build.clean
build.clean:
	@rm -f $(CURDIR)/$(TARGET_NAME)