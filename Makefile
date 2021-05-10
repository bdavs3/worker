windows:
	GOOS=windows go build -o bin/worker cli/cli.go
	GOOS=windows go build -o bin/server server/server.go

mac:
	GOOS=darwin go build -o bin/worker cli/cli.go
	GOOS=darwin go build -o bin/server server/server.go

linux:
	GOOS=linux go build -o bin/worker cli/cli.go
	GOOS=linux go build -o bin/server server/server.go