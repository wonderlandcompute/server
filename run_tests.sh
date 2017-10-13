#!/bin/bash

if [ -z "$DISNEYLAND_TEST_DB" ]; then
    echo "Environment variable \$DISNEYLAND_TEST_DB needs to be set!"
    exit 1
fi

migrate -database $DISNEYLAND_TEST_DB -path migrations down
migrate -database $DISNEYLAND_TEST_DB -path migrations up

cd disneyland; go test -v .; cd ..
