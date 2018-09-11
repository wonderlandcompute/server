image:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/wonderland-server wonderland_server.go
	CGO_ENABLED=0 GOOS=linux go build -tags 'postgres' -a -installsuffix cgo -o build/migrate github.com/mattes/migrate/cli

	docker build -t wondercompute -f Dockerfile.scratch .

build_and_run:	image
	./build/wonderland-server

run:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/wonderland-server wonderland_server.go
	./build/wonderland-server
