.PHONY: build doc local

default: build

build: asapp
	GOOS=linux GOARCH=amd64 go build -o asapp

doc:
	godoc -http=:8765

local:
	go build -o asapp
	./asapp -c=development.json serve
