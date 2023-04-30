EXEC = docuowl

ifeq ($(OS),Windows_NT)
	EXEC := $(EXEC).exe
endif

all:
	go run ./cmd/static-generator/main.go
	go build -o ./$(EXEC) ./cmd/docuowl/main.go
