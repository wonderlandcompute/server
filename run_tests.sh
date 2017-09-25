#!/bin/bash

if [ -z "$OPTIMUS_TEST_DB" ]; then
    echo "Environment variable \$OPTIMUS_TEST_DB needs to be set!"
    exit 1
fi

migrate -path migrations -url $OPTIMUS_TEST_DB reset
cd optimus; go test -v .; cd ..