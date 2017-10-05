FROM golang:latest

# Migration stuff
RUN go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
RUN go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli

# Dep
RUN go get -u github.com/golang/dep/cmd/dep


# Package install
ADD . /go/src/gitlab.com/lambda-hse/optimus
RUN cd /go/src/gitlab.com/lambda-hse/optimus; dep ensure
RUN go install gitlab.com/lambda-hse/optimus


# Startup
# ENTRYPOINT /bin/bash
EXPOSE 50051
