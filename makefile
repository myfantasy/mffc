build:
	CGO_ENABLED=0 go build -o bin/app
run:
	bin/app
br: build run


t: 
	go test fod/*.go
	go test ./compress/
	go test ./compress/ -bench=.