.DEFAULT_GOAL := build

.PHONY:fmt vet build

fmt:
	go fmt look.go

vet: fmt
	go vet look.go

build: vet
	go build look.go

clean: look
	/bin/rm look

