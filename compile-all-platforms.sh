#!/bin/bash

mkdir -p .bin/

go get

for os in windows linux darwin; do
  for arch in amd64 386 arm arm64; do
    target="$os:$arch"
    echo "Compiling $target"
    export GOOS=$(echo $target | cut -d: -f1) GOARCH=$(echo $target | cut -d: -f2)
    OUT=.bin/$(basename $(echo $PWD))_${GOOS}_${GOARCH}
    if [ $GOOS == "windows" ]
    then
      OUT="$OUT.exe"
    fi
    bash -c "go build -ldflags '-w -X main.buildTime=$(date +'%b-%d-%Y-%H:%M:%S') -X main.commit=$(git describe --always)' -v -o $OUT ."
  done
done
