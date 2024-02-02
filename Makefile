build.linux:
	CGO_ENABLED=1 go build -o bin/pm main.go