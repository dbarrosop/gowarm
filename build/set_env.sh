#!/bin/bash
if [ ${1} == "raspberrypi" ]; then
    echo export GOOS=linux
    echo export GOARCH=arm
    echo export CGO_ENABLED=1
    echo unset GOFLAGS
    echo unset GOROOT
else
    echo export GOOS=$(tinygo info -target ${1} | grep GOOS | awk '{print $2}')
    echo export GOARCH=$(tinygo info -target ${1} | grep GOARCH | awk '{print $2}')
    echo export GOFLAGS=-tags=$(tinygo info -target ${1} | grep tags | sed -E "s/build tags:[[:space:]]+//" | sed -E "s/ /,/g")
    echo export GOROOT=$(tinygo info -target ${1} | grep GOROOT | awk '{print $3}')
    echo export CGO_ENABLED=1
fi
