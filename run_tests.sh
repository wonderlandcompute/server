#!/bin/bash

if [ -z "$WONDERLAND_TEST_DB" ]; then
    echo "Environment variable \$WONDERLAND_TEST_DB needs to be set!"
    exit 1
fi

migrate -database $WONDERLAND_TEST_DB -path migrations down
migrate -database $WONDERLAND_TEST_DB -path migrations up

cd wonderland; go test -v .; cd ..
