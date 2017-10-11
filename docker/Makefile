image:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o optimus-server ../optimus_server.go
	docker build -t optimus -f Dockerfile.scratch .