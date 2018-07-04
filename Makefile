image:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/wonderland-server wonder_server.go
	CGO_ENABLED=0 GOOS=linux go build -tags 'postgres' -a -installsuffix cgo -o build/migrate github.com/mattes/migrate/cli

	docker build -t wonderland -f Dockerfile.scratch .
