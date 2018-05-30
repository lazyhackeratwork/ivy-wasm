#!/bin/sh


if [ -x ivy.wasm ]; then
    rm ivy.wasm
fi

GOROOT=$HOME/go-wasm
GOOS=js GOARCH=wasm $HOME/go-wasm/bin/go build -o ivy.wasm ivy.go
