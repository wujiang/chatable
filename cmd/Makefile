.PHONY: build doc local

default: build

build: chatable
	GOOS=linux GOARCH=amd64 go build -o chatable

doc:
	godoc -http=:8765

local:
	go build -o chatable
	./chatable -c=development.json serve
