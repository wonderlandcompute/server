#!/bin/bash

if [ -z "$OPTIMUS_TEST_DB" ]; then
    echo "Environment variable \$OPTIMUS_TEST_DB needs to be set!"
    exit 1
fi

migrate -database $OPTIMUS_TEST_DB -path migrations down
migrate -database $OPTIMUS_TEST_DB -path migrations up

cd optimus; go test -v .; cd ..
