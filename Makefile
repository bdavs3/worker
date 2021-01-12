build:
	GOOS=linux go build -o bin/worker cli/cli.go
	GOOS=linux go build -o bin/server server/server.go