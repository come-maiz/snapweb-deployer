#!/bin/sh

set -eu

if which goctest >/dev/null; then
    goctest="goctest"
else
    goctest="go test"
fi

STATIC=""
UNIT=""

case "${1:-all}" in
    all)
        STATIC="yes"
        UNIT="yes"
        ;;
    --static)
        STATIC="yes"
        ;;
    --unit)
        UNIT="yes"
        ;;
    *)
        echo "Wrong flag ${1}. To run a single suite use --static or --unit."
        exit 1
esac

# Append the coverage profile of a package to the project coverage.
append_go_coverage() {
    local profile="$1"
    if [ -f $profile ]; then
        cat $profile | grep -v "mode: set" >> .coverage-go/coverage.out
        rm $profile
    fi
}

if [ ! -z "$STATIC" ]; then
    # Run static tests.

    echo Checking formatting
    fmt=$(gofmt -l .)

    if [ -n "$fmt" ]; then
        echo "Formatting wrong in following files"
        echo $fmt
        exit 1
    fi

    # go vet
    echo Running vet
    go vet ./...

    echo Install golint
    go get github.com/golang/lint/golint
    export PATH=$PATH:$GOPATH/bin

    echo Running lint
    lint=$(golint ./...)
    if [ -n "$lint" ]; then
        echo "Lint complains:"
        echo $lint
        exit 1
    fi

fi

if [ ! -z "$UNIT" ]; then
    echo Building
    go build -v github.com/elopio/snapweb-deployer/...

    # Prepare the coverage output profile.
    rm -rf .coverage-go
    mkdir .coverage-go
    echo "mode: set" > .coverage-go/coverage.out

    # tests
    echo Running tests from $(pwd)
    for pkg in $(go list ./...); do
        $goctest -v -coverprofile=.coverage-go/profile.out $pkg
        append_go_coverage .coverage-go/profile.out
    done
fi

echo "All good, what could possibly go wrong"
